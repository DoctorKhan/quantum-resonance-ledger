package rtt

import (
	"testing"
	// "time" // Import time if needed for timestamp tests later
)

// TestPropensityField verifies the initialization and basic properties of PropensityField.
func TestPropensityField(t *testing.T) {
	// t.Skip("TDD Step 1: Test for PropensityField now implemented.") // Unskipped
	assetID := "ASSET_PROP"
	pf := PropensityField{
		AssetID: assetID,
		Density: make(map[string]float64), // Explicitly initialize map
	}

	if pf.AssetID != assetID {
		t.Errorf("Expected AssetID to be '%s', got '%s'", assetID, pf.AssetID)
	}
	if pf.Density == nil {
		t.Fatal("Density map was not initialized")
	}
	if len(pf.Density) != 0 {
		t.Errorf("Expected Density map to be empty, but got size %d", len(pf.Density))
	}

	// Test adding a density value
	pr := PriceRange{Min: 1.0, Max: 2.0}
	key := pr.Key()
	pf.Density[key] = 0.5
	if pf.Density[key] != 0.5 {
		t.Errorf("Expected density for key '%s' to be 0.5, got %f", key, pf.Density[key])
	}
	if len(pf.Density) != 1 {
		t.Errorf("Expected Density map size to be 1 after adding, but got size %d", len(pf.Density))
	}
	// TODO: Add tests for serialization if needed later.
}

// TestLocalNodeStateRTT verifies the initialization and structure of LocalNodeStateRTT.
func TestLocalNodeStateRTT(t *testing.T) {
	// t.Skip("TDD Step 1: Test for LocalNodeStateRTT now implemented.") // Unskipped
	state := InitializeRTTState() // Use the initializer

	if state == nil {
		t.Fatal("InitializeRTTState returned nil, cannot test LocalNodeStateRTT structure")
	}

	// Check that the core maps are initialized (not nil)
	if state.BuyFields == nil {
		t.Error("LocalNodeStateRTT.BuyFields map was not initialized")
	}
	if state.SellFields == nil {
		t.Error("LocalNodeStateRTT.SellFields map was not initialized")
	}

	// Check initial map sizes (should be empty)
	if len(state.BuyFields) != 0 {
		t.Errorf("Expected BuyFields map to be empty, but got size %d", len(state.BuyFields))
	}
	if len(state.SellFields) != 0 {
		t.Errorf("Expected SellFields map to be empty, but got size %d", len(state.SellFields))
	}

	// TODO: Add checks for other fields (CUTs, Neighbors, Q) once they are added to the struct and initializer.
}

// TestSettlementRecord verifies the creation and properties of SettlementRecord.
func TestSettlementRecord(t *testing.T) {
	// t.Skip("TDD Step 1: Test for SettlementRecord now implemented.") // Unskipped
	recordID := "settle-123"
	assetID := "ASSET_SR"
	amount := 100.50
	priceRangeStr := "10.00-11.00" // Using string as per current struct definition
	nodeID := "NodeSettle"
	timestamp := int64(1678886400000000000) // Example timestamp

	sr := SettlementRecord{
		RecordID:   recordID,
		AssetID:    assetID,
		Amount:     amount,
		PriceRange: priceRangeStr,
		Timestamp:  timestamp,
		NodeID:     nodeID,
		// TODO: Add Involved CUTs once defined
	}

	if sr.RecordID != recordID {
		t.Errorf("Expected RecordID to be '%s', got '%s'", recordID, sr.RecordID)
	}
	if sr.AssetID != assetID {
		t.Errorf("Expected AssetID to be '%s', got '%s'", assetID, sr.AssetID)
	}
	if sr.Amount != amount {
		t.Errorf("Expected Amount to be %f, got %f", amount, sr.Amount)
	}
	if sr.PriceRange != priceRangeStr {
		t.Errorf("Expected PriceRange to be '%s', got '%s'", priceRangeStr, sr.PriceRange)
	}
	if sr.Timestamp != timestamp {
		t.Errorf("Expected Timestamp to be %d, got %d", timestamp, sr.Timestamp)
	}
	if sr.NodeID != nodeID {
		t.Errorf("Expected NodeID to be '%s', got '%s'", nodeID, sr.NodeID)
	}
	// TODO: Test serialization/deserialization if needed later.
}

// TestPropagationPacket verifies the creation and properties of PropagationPacket.
func TestPropagationPacket(t *testing.T) {
	// t.Skip("TDD Step 1: Test for PropagationPacket now implemented.") // Unskipped
	packetID := "prop-pkt-456"
	sourceNode := "NodeProp"
	timestamp := int64(1678887000000000000) // Example timestamp

	// Create sample updates and reports
	fieldUpdate1 := FieldUpdate{
		AssetID:       "ASSET_P1",
		PriceRange:    PriceRange{Min: 10, Max: 11},
		IsBuyField:    true,
		DensityChange: 0.1,
	}
	settlementReport1 := SettlementRecord{
		RecordID:   "settle-prop-1",
		AssetID:    "ASSET_P2",
		Amount:     50.0,
		PriceRange: "12.00-13.00",
		Timestamp:  timestamp - 1000000000, // Earlier timestamp
		NodeID:     sourceNode,
	}

	pp := PropagationPacket{
		PacketID:          packetID,
		SourceNode:        sourceNode,
		Timestamp:         timestamp,
		FieldUpdates:      []FieldUpdate{fieldUpdate1},           // Initialize with one update
		SettlementReports: []SettlementRecord{settlementReport1}, // Initialize with one report
	}

	if pp.PacketID != packetID {
		t.Errorf("Expected PacketID to be '%s', got '%s'", packetID, pp.PacketID)
	}
	if pp.SourceNode != sourceNode {
		t.Errorf("Expected SourceNode to be '%s', got '%s'", sourceNode, pp.SourceNode)
	}
	if pp.Timestamp != timestamp {
		t.Errorf("Expected Timestamp to be %d, got %d", timestamp, pp.Timestamp)
	}

	// Check slices are initialized and contain the expected data
	if pp.FieldUpdates == nil {
		t.Fatal("FieldUpdates slice was not initialized")
	}
	if len(pp.FieldUpdates) != 1 {
		t.Fatalf("Expected FieldUpdates slice to have length 1, got %d", len(pp.FieldUpdates))
	}
	if pp.FieldUpdates[0] != fieldUpdate1 {
		t.Errorf("FieldUpdates content mismatch: expected %+v, got %+v", fieldUpdate1, pp.FieldUpdates[0])
	}

	if pp.SettlementReports == nil {
		t.Fatal("SettlementReports slice was not initialized")
	}
	if len(pp.SettlementReports) != 1 {
		t.Fatalf("Expected SettlementReports slice to have length 1, got %d", len(pp.SettlementReports))
	}
	if pp.SettlementReports[0] != settlementReport1 {
		t.Errorf("SettlementReports content mismatch: expected %+v, got %+v", settlementReport1, pp.SettlementReports[0])
	}

	// Test initialization with empty slices
	ppEmpty := PropagationPacket{
		PacketID:          "prop-empty",
		SourceNode:        "NodeEmpty",
		Timestamp:         timestamp,
		FieldUpdates:      make([]FieldUpdate, 0),      // Explicitly initialize empty
		SettlementReports: make([]SettlementRecord, 0), // Explicitly initialize empty
	}
	if ppEmpty.FieldUpdates == nil {
		t.Error("Empty FieldUpdates slice should be initialized, not nil")
	}
	if len(ppEmpty.FieldUpdates) != 0 {
		t.Errorf("Expected empty FieldUpdates slice length to be 0, got %d", len(ppEmpty.FieldUpdates))
	}
	if ppEmpty.SettlementReports == nil {
		t.Error("Empty SettlementReports slice should be initialized, not nil")
	}
	if len(ppEmpty.SettlementReports) != 0 {
		t.Errorf("Expected empty SettlementReports slice length to be 0, got %d", len(ppEmpty.SettlementReports))
	}

	// TODO: Test serialization/deserialization if needed later.
}

// TestPriceRange verifies the basic properties of PriceRange.
func TestPriceRange(t *testing.T) {
	// t.Skip("TDD Step 1: Test for PriceRange now implemented.") // Unskipped
	pr := PriceRange{Min: 10.50, Max: 11.00}
	if pr.Min != 10.50 {
		t.Errorf("Expected Min to be 10.50, got %f", pr.Min)
	}
	if pr.Max != 11.00 {
		t.Errorf("Expected Max to be 11.00, got %f", pr.Max)
	}

	expectedKey := "10.50-11.00"
	actualKey := pr.Key()
	if actualKey != expectedKey {
		t.Errorf("Expected Key() to return '%s', got '%s'", expectedKey, actualKey)
	}

	// Test Key() with different precision
	pr2 := PriceRange{Min: 9.99, Max: 10.01}
	expectedKey2 := "9.99-10.01"
	actualKey2 := pr2.Key()
	if actualKey2 != expectedKey2 {
		t.Errorf("Expected Key() to return '%s', got '%s'", expectedKey2, actualKey2)
	}

	// TODO: Consider adding validation test (Min <= Max) if a validation method is added.
}

// TestFieldUpdate verifies the basic properties of FieldUpdate.
func TestFieldUpdate(t *testing.T) {
	// t.Skip("TDD Step 1: Test for FieldUpdate now implemented.") // Unskipped
	assetID := "ASSET_X"
	priceRange := PriceRange{Min: 50.0, Max: 51.0}
	densityChange := 0.15

	fu := FieldUpdate{
		AssetID:       assetID,
		PriceRange:    priceRange,
		IsBuyField:    true,
		DensityChange: densityChange,
	}

	if fu.AssetID != assetID {
		t.Errorf("Expected AssetID to be '%s', got '%s'", assetID, fu.AssetID)
	}
	if fu.PriceRange != priceRange {
		t.Errorf("Expected PriceRange to be %+v, got %+v", priceRange, fu.PriceRange)
	}
	if !fu.IsBuyField {
		t.Error("Expected IsBuyField to be true, got false")
	}
	if fu.DensityChange != densityChange {
		t.Errorf("Expected DensityChange to be %f, got %f", densityChange, fu.DensityChange)
	}

	// Test with IsBuyField = false
	fuSell := FieldUpdate{
		AssetID:       assetID,
		PriceRange:    priceRange,
		IsBuyField:    false,
		DensityChange: -0.05,
	}
	if fuSell.IsBuyField {
		t.Error("Expected IsBuyField to be false, got true")
	}
	if fuSell.DensityChange != -0.05 {
		t.Errorf("Expected DensityChange to be -0.05, got %f", fuSell.DensityChange)
	}
}

// TestInitializeRTTState verifies the RTT state initialization function.
func TestInitializeRTTState(t *testing.T) {
	// t.Skip("TDD Step 1: Test for InitializeRTTState now implemented.") // Unskipped
	state := InitializeRTTState()
	if state == nil {
		t.Fatal("InitializeRTTState returned nil")
	}
	if state.BuyFields == nil {
		t.Error("InitializeRTTState did not initialize BuyFields map")
	}
	if state.SellFields == nil {
		t.Error("InitializeRTTState did not initialize SellFields map")
	}
	// TODO: Add checks for other initialized fields (CUTs map, Neighbor map, Q map) once defined.
}
