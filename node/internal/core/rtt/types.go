package rtt

import "fmt" // Added for PriceRange.Key()

// Placeholder for RTT core data structures.
// Implementation will follow TDD steps outlined in docs/tdd-rtt-go.md.

// PropensityField represents the probability density of trading intent
// over asset and price ranges. The exact representation (e.g., histogram,
// spline coefficients, discrete map) needs further definition based on
// the required mathematical operations (overlap integral, perturbation).
type PropensityField struct {
	// Example: Using a map for discrete price buckets
	// Density map[PriceBucket]float64 // PriceBucket needs definition
	// Or potentially a more complex structure for continuous fields.
	// For now, keep it minimal.
	AssetID string
	// Density represents the probability density, keyed by a PriceRange representation.
	// Using string key for simplicity initially. Map key could be PriceRange itself if comparable.
	Density map[string]float64
	// TODO: Revisit this representation based on required math operations.
}

// LocalNodeStateRTT holds the RTT-specific state relevant for local calculations.
type LocalNodeStateRTT struct {
	// Own propensity fields, keyed by AssetID
	BuyFields  map[string]*PropensityField
	SellFields map[string]*PropensityField

	// TODO: Representation of relevant local CUTs (e.g., map[CUT_ID]CUT_Status)
	// TODO: Snapshot of recently received neighbor data (e.g., map[NeighborID]NeighborSnapshot)
	// TODO: Local Quantity Imbalance (Q) for various assets map[AssetID]float64
}

// SettlementRecord captures the essential details of a probabilistic local settlement event.
type SettlementRecord struct {
	RecordID   string // Unique identifier for this settlement event
	AssetID    string
	Amount     float64 // Consider using a high-precision decimal type for financial calculations
	PriceRange string  // TODO: Define a specific PriceRange type/struct
	// TODO: List or map of involved CUT identifiers (e.g., []CUT_ID)
	Timestamp int64  // Nanoseconds since epoch, for ordering and propagation timing
	NodeID    string // ID of the node where settlement occurred
}

// PropagationPacket contains the information nodes exchange to update RTT state.
type PropagationPacket struct {
	PacketID   string // Unique ID for the packet
	SourceNode string
	Timestamp  int64

	// Updates to propensity fields from the source node
	FieldUpdates []FieldUpdate // TODO: Define FieldUpdate struct (e.g., AssetID, PriceRange, DeltaDensity)

	// Reports of settlements that occurred at the source node
	SettlementReports []SettlementRecord

	// TODO: Potentially include Q imbalance updates or other necessary sync data.
}

// --- Supporting Type Definitions (Placeholders) ---

// PriceRange might define a specific range for propensity fields or settlements.
type PriceRange struct {
	Min float64
	Max float64
}

// Key generates a simple string representation for use as a map key.
// TODO: Consider a more robust hashing or canonical representation if needed.
func (pr PriceRange) Key() string {
	// Note: Formatting needs care to ensure uniqueness and consistency.
	// Using simple format for now. Consider fixed precision.
	return fmt.Sprintf("%.2f-%.2f", pr.Min, pr.Max)
}

// FieldUpdate could represent a change in a specific part of a propensity field.
type FieldUpdate struct {
	AssetID       string
	PriceRange    PriceRange // Or PriceBucket
	IsBuyField    bool
	DensityChange float64
}

// CUT_ID placeholder type
// type CUT_ID string

// InitializeRTTState creates a new, empty RTT state for a node.
func InitializeRTTState() *LocalNodeStateRTT {
	return &LocalNodeStateRTT{
		BuyFields:  make(map[string]*PropensityField),
		SellFields: make(map[string]*PropensityField),
		// Initialize other maps/slices
	}
}
