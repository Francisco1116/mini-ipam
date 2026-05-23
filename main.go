package main

import (
	"net/http"
	"net/netip"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AllocateSubnetRequest struct {
	Name   string `json:"name" binding:"required"`
	Prefix int    `json:"prefix" binding:"required"`
}

type Subnet struct {
	ID        string `json:"id"`
	NetworkID string `json:"network_id"`
	Name      string `json:"name"`
	CIDR      string `json:"cidr"`
}

type Network struct {
	ID string `json:"id"`
	Name string `json:"name"`
	CIDR string `json:"cidr"`
	Subnets []Subnet `json:"subnets"`
}

type CreateNetworkRequest struct {
	Name string `json:"name" binding:"required"`
	CIDR string `json:"cidr" binding:"required"`
}

func ipToUint32(ip netip.Addr) uint32 {
	b := ip.As4()
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func uint32ToIP(val uint32) netip.Addr {
	return netip.AddrFrom4([4]byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)})
}

var networkDB = make(map[string]Network)

func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
    v1.POST("/networks", createNetwork)
		v1.GET("/networks/:id", getNetwork)

		v1.POST("/networks/:id/subnets", allocateSubnet)
	}

	r.Run(":8080")
}

func createNetwork(c *gin.Context) {
	var req CreateNetworkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing fields"})
		return
	}

	  prefix, err := netip.ParsePrefix(req.CIDR)
	  if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CIDR format"})
		return
	  }

	if prefix.Masked() != prefix {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CIDR must be a valid network address (e.g., 10.0.0.0/16)"})
		return
	}

	id := uuid.New().String()
	newNet := Network{
		ID:      id,
		Name:    req.Name,
		CIDR:    prefix.String(),
		Subnets: []Subnet{},
	}

	networkDB[id] = newNet

	c.JSON(http.StatusCreated, newNet)
}

func getNetwork(c *gin.Context) {
	id := c.Param("id")

	net, exists := networkDB[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Network not found"})
		return
	}

	c.JSON(http.StatusOK, net)
}

func allocateSubnet(c *gin.Context) {
	id := c.Param("id")
	var req AllocateSubnetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	netRecord, exists := networkDB[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Network not found"})
		return
	}

	parentPrefix, _ := netip.ParsePrefix(netRecord.CIDR)
	reqPrefixSize := req.Prefix

	if reqPrefixSize < parentPrefix.Bits() || reqPrefixSize > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subnet prefix size"})
		return
	}

	step := uint32(1 << (32 - reqPrefixSize)) 
	startIP := ipToUint32(parentPrefix.Addr())
	maxIP := startIP + uint32(1<<(32-parentPrefix.Bits())) - 1

	var allocated *Subnet

	for current := startIP; current+step-1 <= maxIP; current += step {
		candidateIP := uint32ToIP(current)
		candidatePrefix := netip.PrefixFrom(candidateIP, reqPrefixSize)

		overlap := false
		for _, ex := range netRecord.Subnets {
			exPrefix, _ := netip.ParsePrefix(ex.CIDR)
			if candidatePrefix.Overlaps(exPrefix) { 
				overlap = true
				break
			}
		}

		if !overlap {
			newSubnet := Subnet{
				ID:        uuid.New().String(),
				NetworkID: id,
				Name:      req.Name,
				CIDR:      candidatePrefix.String(),
			}
			netRecord.Subnets = append(netRecord.Subnets, newSubnet)
			networkDB[id] = netRecord
			allocated = &newSubnet
			break
		}
	}

	if allocated == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "No available IP space for the requested subnet size"})
		return
	}

	c.JSON(http.StatusCreated, allocated)
}