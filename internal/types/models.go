package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// CacheControl represents the cache control configuration
type CacheControl struct {
	Type string `json:"type"`          // "ephemeral"
	TTL  string `json:"ttl,omitempty"` // "5m" or "1h", omitted if empty
}

// UnmarshalJSON implements custom unmarshaling for CacheControl to handle empty TTL
func (cc *CacheControl) UnmarshalJSON(data []byte) error {
	type Alias CacheControl
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(cc),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Set default TTL if empty
	if cc.TTL == "" {
		cc.TTL = "5m"
	}
	return nil
}

// ContentBlock represents a content block in a message
type ContentBlock struct {
	Type         string        `json:"type"`
	Text         string        `json:"text,omitempty"`
	Data         string        `json:"data,omitempty"`
	Source       *ImageSource  `json:"source,omitempty"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`

	// For tool_use blocks (assistant messages)
	ID    string      `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Input interface{} `json:"input,omitempty"`

	// For tool_result blocks (user messages)
	ToolUseID string      `json:"tool_use_id,omitempty"`
	Content   interface{} `json:"content,omitempty"`
	IsError   *bool       `json:"is_error,omitempty"`

	// For thinking blocks (assistant messages)
	Thinking  string `json:"thinking,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// ImageSource represents an image source
type ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// UnmarshalJSON implements custom unmarshaling for Message to support both string and array content formats
func (m *Message) UnmarshalJSON(data []byte) error {
	// Create a temporary struct with Content as json.RawMessage to inspect it first
	type Alias Message
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Check if content is a string (shorthand format)
	var contentStr string
	if err := json.Unmarshal(aux.Content, &contentStr); err == nil {
		// It's a string, convert to ContentBlock array
		m.Content = []ContentBlock{
			{
				Type: "text",
				Text: contentStr,
			},
		}
		return nil
	}

	// It's an array, unmarshal normally
	var contentBlocks []ContentBlock
	if err := json.Unmarshal(aux.Content, &contentBlocks); err != nil {
		return err
	}
	m.Content = contentBlocks

	return nil
}

// ToolDefinition represents a tool definition
type ToolDefinition struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	InputSchema  interface{}   `json:"input_schema"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

// AnthropicRequest represents the complete request to Anthropic API
type AnthropicRequest struct {
	Model         string           `json:"model"`
	MaxTokens     int              `json:"max_tokens"`
	Messages      []Message        `json:"messages"`
	System        string           `json:"system,omitempty"`
	SystemBlocks  []ContentBlock   `json:"-"` // Handle system blocks through custom parsing
	Tools         []ToolDefinition `json:"tools,omitempty"`
	Temperature   *float64         `json:"temperature,omitempty"`
	TopP          *float64         `json:"top_p,omitempty"`
	TopK          *int             `json:"top_k,omitempty"`
	Stream        *bool            `json:"stream,omitempty"`
	StopSequences []string         `json:"stop_sequences,omitempty"`
}

// MarshalJSON implements custom marshaling for AnthropicRequest
// When SystemBlocks is populated, it serializes as the "system" field instead of the System string
func (r *AnthropicRequest) MarshalJSON() ([]byte, error) {
	// Create an alias type to avoid infinite recursion
	type Alias AnthropicRequest

	// If SystemBlocks is populated, we need to serialize it as "system" array
	if len(r.SystemBlocks) > 0 {
		// Create an anonymous struct with all fields explicitly set
		return json.Marshal(&struct {
			Model         string           `json:"model"`
			MaxTokens     int              `json:"max_tokens"`
			Messages      []Message        `json:"messages"`
			System        interface{}      `json:"system,omitempty"` // Use interface{} to allow array
			Tools         []ToolDefinition `json:"tools,omitempty"`
			Temperature   *float64         `json:"temperature,omitempty"`
			TopP          *float64         `json:"top_p,omitempty"`
			TopK          *int             `json:"top_k,omitempty"`
			Stream        *bool            `json:"stream,omitempty"`
			StopSequences []string         `json:"stop_sequences,omitempty"`
		}{
			Model:         r.Model,
			MaxTokens:     r.MaxTokens,
			Messages:      r.Messages,
			System:        r.SystemBlocks, // Serialize SystemBlocks as "system"
			Tools:         r.Tools,
			Temperature:   r.Temperature,
			TopP:          r.TopP,
			TopK:          r.TopK,
			Stream:        r.Stream,
			StopSequences: r.StopSequences,
		})
	}

	// Otherwise, use normal marshaling (System string field)
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom unmarshaling for AnthropicRequest
// Handles both string and array formats for the "system" field
func (r *AnthropicRequest) UnmarshalJSON(data []byte) error {
	// Create an alias type to avoid infinite recursion
	type Alias AnthropicRequest

	// First, try to unmarshal with a temporary struct that has system as RawMessage
	aux := &struct {
		*Alias
		System json.RawMessage `json:"system,omitempty"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// If there's a system field, check if it's a string or array
	if len(aux.System) > 0 {
		// Try to unmarshal as string first
		var systemStr string
		if err := json.Unmarshal(aux.System, &systemStr); err == nil {
			r.System = systemStr
			return nil
		}

		// If that fails, try to unmarshal as array of content blocks
		var systemBlocks []ContentBlock
		if err := json.Unmarshal(aux.System, &systemBlocks); err == nil {
			r.SystemBlocks = systemBlocks
			r.System = "" // Clear the string field
			return nil
		}

		// If both fail, return error
		return fmt.Errorf("system field must be either a string or an array of content blocks")
	}

	return nil
}

// AnthropicResponse represents the response from Anthropic API
type AnthropicResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        Usage          `json:"usage"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
}

// CacheBreakpoint represents a cache breakpoint decision
type CacheBreakpoint struct {
	Position    string    `json:"position"`      // "system", "tools", "message_0_block_1"
	Tokens      int       `json:"tokens"`        // Number of tokens cached
	TTL         string    `json:"ttl,omitempty"` // "5m" or "1h", omitted if empty
	Type        string    `json:"type"`          // "system", "tools", "content"
	WritePrice  float64   `json:"write_price"`   // Cost to write this cache
	ReadSavings float64   `json:"read_savings"`  // Savings per read
	Timestamp   time.Time `json:"timestamp"`
}

// ROIMetrics represents return on investment calculations
type ROIMetrics struct {
	BaseInputCost        float64 `json:"base_input_cost"`         // Original cost without caching
	CacheWriteCost       float64 `json:"cache_write_cost"`        // Additional cost for cache writes
	CacheReadCost        float64 `json:"cache_read_cost"`         // Cost for cache reads (subsequent requests)
	FirstRequestCost     float64 `json:"first_request_cost"`      // Total cost including cache writes
	SubsequentSavings    float64 `json:"subsequent_savings"`      // Savings per subsequent request
	BreakEvenRequests    int     `json:"break_even_requests"`     // Number of requests to break even
	SavingsAt10Requests  float64 `json:"savings_at_10_requests"`  // Total savings after 10 requests
	SavingsAt100Requests float64 `json:"savings_at_100_requests"` // Total savings after 100 requests
	PercentSavings       float64 `json:"percent_savings"`         // Percentage savings at scale
}

// CacheMetadata represents metadata about caching decisions
type CacheMetadata struct {
	CacheInjected bool              `json:"cache_injected"`
	TotalTokens   int               `json:"total_tokens"`
	CachedTokens  int               `json:"cached_tokens"`
	CacheRatio    float64           `json:"cache_ratio"` // Percentage of tokens cached
	Breakpoints   []CacheBreakpoint `json:"breakpoints"`
	ROI           ROIMetrics        `json:"roi"`
	Strategy      string            `json:"strategy"` // "aggressive", "moderate", "conservative"
	Model         string            `json:"model"`
	Timestamp     time.Time         `json:"timestamp"`
}

// CacheStrategy represents different caching strategies
type CacheStrategy string

const (
	StrategyConservative   CacheStrategy = "conservative"
	StrategyModerate       CacheStrategy = "moderate"
	StrategyAggressive     CacheStrategy = "aggressive"
	StrategyAutoAggressive CacheStrategy = "auto_aggressive"
)

// CacheCandidate represents a potential cache breakpoint
// Moved from cache/injector.go for shared use across packages
type CacheCandidate struct {
	Position    string      // "tools", "system", "message_X_block_1", etc.
	Tokens      int         // Token count
	ContentType string      // "tools", "system", "content"
	TTL         string      // "5m" or "1h"
	ROIScore    float64     // ROI score for prioritization
	WriteCost   float64     // Cost to write cache
	ReadSavings float64     // Savings per read
	BreakEven   int         // Requests to break even
	Content     interface{} // Reference to the actual content
}

// BreakpointSelector contains strategy-specific breakpoint selection logic
type BreakpointSelector struct {
	// ShouldAcceptExisting returns true if we should reuse existing breakpoints
	// Returns false to recalculate breakpoints from scratch
	// Now receives strategy and request context for intelligent decisions
	ShouldAcceptExisting func(existing []CacheCandidate, strategy CacheStrategy, req *AnthropicRequest) bool

	// SelectBreakpoints filters and selects breakpoints according to strategy rules
	// Input candidates are sorted: [tools, system, content1, content2, ...]
	// Now receives strategy and request context for intelligent decisions
	// MaxBreakpoints can be obtained from GetStrategyConfig(strategy).MaxBreakpoints
	SelectBreakpoints func(candidates []CacheCandidate, strategy CacheStrategy, req *AnthropicRequest) []CacheCandidate
}

// StrategyConfig represents configuration for each strategy
type StrategyConfig struct {
	MaxBreakpoints      int                `json:"max_breakpoints"`
	MinTokensMultiplier float64            `json:"min_tokens_multiplier"` // Multiplier for base minimum tokens
	SystemTTL           string             `json:"system_ttl"`
	ToolsTTL            string             `json:"tools_ttl"`
	ContentTTL          string             `json:"content_ttl"`
	Priority            []string           `json:"priority"` // Order of content types to prioritize
	Selector            BreakpointSelector // Strategy-specific selection logic
}

// GetStrategyConfig returns configuration for a given strategy
func GetStrategyConfig(strategy CacheStrategy) StrategyConfig {
	// Standard selector: always recalculate, sort by ROI
	standardShouldAcceptExisting := func(existing []CacheCandidate, strategy CacheStrategy, req *AnthropicRequest) bool {
		return false // Always recalculate for conservative/moderate/aggressive
	}

	standardSelectBreakpoints := func(candidates []CacheCandidate, strategy CacheStrategy, req *AnthropicRequest) []CacheCandidate {
		// Get maxBreakpoints from strategy config
		maxBreakpoints := GetStrategyConfig(strategy).MaxBreakpoints

		// Sort by ROI score (descending)
		for i := 0; i < len(candidates); i++ {
			for j := i + 1; j < len(candidates); j++ {
				if candidates[j].ROIScore > candidates[i].ROIScore {
					candidates[i], candidates[j] = candidates[j], candidates[i]
				}
			}
		}
		// Limit to maxBreakpoints
		if len(candidates) > maxBreakpoints {
			return candidates[:maxBreakpoints]
		}
		return candidates
	}

	// auto_aggressive selector: conditionally accept existing breakpoints
	autoAggressiveShouldAcceptExisting := func(existing []CacheCandidate, strategy CacheStrategy, req *AnthropicRequest) bool {
		if len(existing) > 4 {
			return false
		}

		// 统计 tools 和 system 断点数量
		toolsAndSystem := 0
		messagesCount := 0
		for _, c := range existing {
			if c.ContentType == "tools" || c.ContentType == "system" {
				toolsAndSystem++
			} else if c.ContentType == "content" {
				messagesCount++
			}
		}

		// 新逻辑：tools和系统断点数<=2 并且 (messages断点数==消息数 或 messages断点数>=2)
		return toolsAndSystem <= 2 && (messagesCount >= 2 || messagesCount == len(req.Messages))
	}

	autoAggressiveSelectBreakpoints := func(candidates []CacheCandidate, strategy CacheStrategy, req *AnthropicRequest) []CacheCandidate {
		// Get maxBreakpoints from strategy config
		maxBreakpoints := GetStrategyConfig(strategy).MaxBreakpoints

		// Input candidates are sorted: [tools, system, content1, content2, ...]
		var result []CacheCandidate
		var toolsBP *CacheCandidate
		var systemBPs []CacheCandidate
		var contentBPs []CacheCandidate

		// Categorize candidates
		for i := range candidates {
			switch candidates[i].ContentType {
			case "tools":
				candidates[i].TTL = "1h" // Upgrade tools to 1h TTL
				toolsBP = &candidates[i]
			case "system":
				candidates[i].TTL = "1h" // Upgrade tools to 1h TTL
				systemBPs = append(systemBPs, candidates[i])
			case "content":
				contentBPs = append(contentBPs, candidates[i])
			}
		}

		// Step 1: Handle tools + system (target <= 2)
		toolsAndSystemCount := 0
		if toolsBP != nil {
			toolsAndSystemCount++
		}
		toolsAndSystemCount += len(systemBPs)

		if toolsAndSystemCount <= 2 {
			// Keep all tools and system
			if toolsBP != nil {
				result = append(result, *toolsBP)
			}
			result = append(result, systemBPs...)
		} else {
			// Prioritize tools, otherwise keep first system
			if toolsBP != nil {
				result = append(result, *toolsBP)
				// Add last system to reach 2 total
				if len(systemBPs) > 0 {
					result = append(result, systemBPs[len(systemBPs)-1])
				}
			} else {
				// No tools, keep first and last system (up to 2)
				if len(systemBPs) > 0 {
					result = append(result, systemBPs[0])
				}
				if len(systemBPs) > 1 {
					result = append(result, systemBPs[len(systemBPs)-1])
				}
			}
		}

		// Step 2: Handle content breakpoints (take last 2, upgrade second-to-last to 1h)
		if len(contentBPs) >= 2 {
			// Take last 2
			lastTwo := contentBPs[len(contentBPs)-2:]
			// Upgrade second-to-last to 1h
			lastTwo[0].TTL = "1h"
			result = append(result, lastTwo...)
		} else if len(contentBPs) == 1 {
			result = append(result, contentBPs[0])
		}

		// Ensure total doesn't exceed maxBreakpoints
		if len(result) > maxBreakpoints {
			result = result[:maxBreakpoints]
		}

		return result
	}

	configs := map[CacheStrategy]StrategyConfig{
		StrategyConservative: {
			MaxBreakpoints:      2,
			MinTokensMultiplier: 2.0, // More strict token requirements
			SystemTTL:           "1h",
			ToolsTTL:            "1h",
			ContentTTL:          "5m",
			Priority:            []string{"system", "tools"},
			Selector: BreakpointSelector{
				ShouldAcceptExisting: standardShouldAcceptExisting,
				SelectBreakpoints:    standardSelectBreakpoints,
			},
		},
		StrategyModerate: {
			MaxBreakpoints:      3,
			MinTokensMultiplier: 1.0, // Standard token requirements
			SystemTTL:           "1h",
			ToolsTTL:            "1h",
			ContentTTL:          "5m",
			Priority:            []string{"system", "tools", "content"},
			Selector: BreakpointSelector{
				ShouldAcceptExisting: standardShouldAcceptExisting,
				SelectBreakpoints:    standardSelectBreakpoints,
			},
		},
		StrategyAggressive: {
			MaxBreakpoints:      4,
			MinTokensMultiplier: 0.8, // More lenient token requirements
			SystemTTL:           "1h",
			ToolsTTL:            "1h",
			ContentTTL:          "5m",
			Priority:            []string{"system", "tools", "content", "large_content"},
			Selector: BreakpointSelector{
				ShouldAcceptExisting: standardShouldAcceptExisting,
				SelectBreakpoints:    autoAggressiveSelectBreakpoints,
			},
		},
		StrategyAutoAggressive: {
			MaxBreakpoints:      4,
			MinTokensMultiplier: 0.8, // More lenient token requirements
			SystemTTL:           "1h",
			ToolsTTL:            "1h",
			ContentTTL:          "5m",
			Priority:            []string{"system", "tools", "content", "large_content"},
			Selector: BreakpointSelector{
				ShouldAcceptExisting: autoAggressiveShouldAcceptExisting,
				SelectBreakpoints:    autoAggressiveSelectBreakpoints,
			},
		},
	}
	return configs[strategy]
}

// ToHeaderValue converts a struct to a compact string for headers
func (cm *CacheMetadata) ToHeaderValue() string {
	data, _ := json.Marshal(cm)
	return string(data)
}

// ToBreakpointsHeader converts breakpoints to a compact header string
func (cm *CacheMetadata) ToBreakpointsHeader() string {
	if len(cm.Breakpoints) == 0 {
		return ""
	}

	result := ""
	for i, bp := range cm.Breakpoints {
		if i > 0 {
			result += ","
		}
		result += bp.Position + ":" + string(rune(bp.Tokens)) + ":" + bp.TTL
	}
	return result
}
