package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ==================== Error Definitions ====================
var (
	// ErrInvalidEnum is returned when an enum value doesn't match any valid constant
	ErrInvalidEnum = errors.New("invalid enum value")

	// ErrInvalidEnumType is returned when scanning a value that isn't string or []byte
	ErrInvalidEnumType = errors.New("invalid enum type")

	// ErrInvalidRecipient is return when notifications must not have either user_id or group_id
	ErrInvalidRecipient = errors.New("Notification must have either user_id or group_id, not both")
)

// scanEnum is generic helper function for safely scaning enum values from database
func scanEnum[T ~string](dest *T, value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*dest = T(v)
	case []byte:
		*dest = T(string(v))
	default:
		return fmt.Errorf("%w: cannot scan %T into enum", ErrInvalidEnumType, value)
	}
	return nil
}

// validateEnum validates an enum value and returns a standardized error
func validateEnum[T ~string](value T, typeName string, isValid func(T) bool) error {
	if value != "" && !isValid(value) {
		return fmt.Errorf("%w: %s=%q", ErrInvalidEnum, typeName, value)
	}
	return nil
}

// ==================== UserRole ====================

// UserRole represents user role enum
type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleLeader UserRole = "leader"
	UserRoleUser   UserRole = "user"
)

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin, UserRoleLeader, UserRoleUser:
		return true
	default:
		return false
	}
}

func (r *UserRole) Scan(value interface{}) error {
	if err := scanEnum(r, value); err != nil {
		return err
	}
	return validateEnum(*r, "UserRole", UserRole.IsValid)
}

func (r UserRole) Value() (driver.Value, error) {
	if err := validateEnum(r, "UserRole", UserRole.IsValid); err != nil {
		return nil, err
	}
	return string(r), nil
}

// ==================== AuthProvider ====================

// AuthProvider represents authentication provider enum
type AuthProvider string

const (
	AuthProviderEmail    AuthProvider = "email"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderFacebook AuthProvider = "facebook"
	AuthProviderApple    AuthProvider = "apple"
)

func (a AuthProvider) IsValid() bool {
	switch a {
	case AuthProviderEmail, AuthProviderGoogle, AuthProviderFacebook, AuthProviderApple:
		return true
	default:
		return false
	}
}

func (a *AuthProvider) Scan(value interface{}) error {
	if err := scanEnum(a, value); err != nil {
		return err
	}
	return validateEnum(*a, "AuthProvider", AuthProvider.IsValid)
}

func (a AuthProvider) Value() (driver.Value, error) {
	if err := validateEnum(a, "AuthProvider", AuthProvider.IsValid); err != nil {
		return nil, err
	}
	return string(a), nil
}

// ==================== SurvivalStatus ====================

// SurvivalStatus represents user survival status enum
type SurvivalStatus string

const (
	SurvivalStatusPreparing SurvivalStatus = "preparing"
	SurvivalStatusSafe      SurvivalStatus = "safe"
	SurvivalStatusDanger    SurvivalStatus = "danger"
	SurvivalStatusMissing   SurvivalStatus = "missing"
	SurvivalStatusRescued   SurvivalStatus = "rescued"
)

func (s SurvivalStatus) IsValid() bool {
	switch s {
	case SurvivalStatusPreparing, SurvivalStatusSafe, SurvivalStatusDanger, SurvivalStatusMissing, SurvivalStatusRescued:
		return true
	default:
		return false
	}
}

func (s *SurvivalStatus) Scan(value interface{}) error {
	if err := scanEnum(s, value); err != nil {
		return err
	}
	return validateEnum(*s, "SurvivalStatus", SurvivalStatus.IsValid)
}

func (s SurvivalStatus) Value() (driver.Value, error) {
	if err := validateEnum(s, "SurvivalStatus", SurvivalStatus.IsValid); err != nil {
		return nil, err
	}
	return string(s), nil
}

// ==================== ChecklistCategory ====================

// ChecklistCategory represents checklist item category enum

type ChecklistCategory string

const (
	ChecklistCategorySupplies      ChecklistCategory = "supplies"
	ChecklistCategoryDocuments     ChecklistCategory = "documents"
	ChecklistCategoryEmergency     ChecklistCategory = "emergency"
	ChecklistCategoryFood          ChecklistCategory = "food"
	ChecklistCategoryWater         ChecklistCategory = "water"
	ChecklistCategoryShelter       ChecklistCategory = "shelter"
	ChecklistCategoryFirstAid      ChecklistCategory = "first_aid"
	ChecklistCategoryTools         ChecklistCategory = "tools"
	ChecklistCategoryCommunication ChecklistCategory = "communication"
	ChecklistCategoryClothing      ChecklistCategory = "clothing"
)

func (c ChecklistCategory) IsValid() bool {
	switch c {
	case ChecklistCategorySupplies, ChecklistCategoryDocuments, ChecklistCategoryEmergency, ChecklistCategoryFood, ChecklistCategoryWater,
		ChecklistCategoryShelter, ChecklistCategoryFirstAid, ChecklistCategoryTools, ChecklistCategoryCommunication, ChecklistCategoryClothing:
		return true
	default:
		return false
	}
}

func (c *ChecklistCategory) Scan(value interface{}) error {
	if err := scanEnum(c, value); err != nil {
		return err
	}
	return validateEnum(*c, "ChecklistCategory", ChecklistCategory.IsValid)
}

func (c ChecklistCategory) Value() (driver.Value, error) {
	if err := validateEnum(c, "ChecklistCategory", ChecklistCategory.IsValid); err != nil {
		return nil, err
	}
	return string(c), nil
}

// ==================== PriorityLevel ====================

// PriorityLevel represents priority level enum

type PriorityLevel string

const (
	PriorityLevelLow      PriorityLevel = "low"
	PriorityLevelMedium   PriorityLevel = "medium"
	PriorityLevelHigh     PriorityLevel = "high"
	PriorityLevelCritical PriorityLevel = "critical"
)

func (p PriorityLevel) IsValid() bool {
	switch p {
	case PriorityLevelLow, PriorityLevelMedium, PriorityLevelHigh, PriorityLevelCritical:
		return true
	default:
		return false
	}
}

func (p *PriorityLevel) Scan(value interface{}) error {
	if err := scanEnum(p, value); err != nil {
		return err
	}
	return validateEnum(*p, "PriorityLevel", PriorityLevel.IsValid)
}

func (p PriorityLevel) Value() (driver.Value, error) {
	if err := validateEnum(p, "PriorityLevel", PriorityLevel.IsValid); err != nil {
		return nil, err
	}

	return string(p), nil
}

// IsUrgent checks if priority is high or critical
func (p PriorityLevel) IsUrgent() bool {
	return p == PriorityLevelHigh || p == PriorityLevelCritical
}

// ==================== GuideCategory ====================

// GuideCategory represents survival guide category enum
type GuideCategory string

const (
	GuideCategoryFirstAid   GuideCategory = "first_aid"
	GuideCategoryShelter    GuideCategory = "shelter"
	GuideCategoryFood       GuideCategory = "food"
	GuideCategoryWater      GuideCategory = "water"
	GuideCategoryNavigation GuideCategory = "navigation"
	GuideCategoryFire       GuideCategory = "fire"
	GuideCategorySignaling  GuideCategory = "signaling"
	GuideCategoryWeather    GuideCategory = "weather"
	GuideCategoryWildlife   GuideCategory = "wildlife"
	GuideCategoryTools      GuideCategory = "tools"
)

func (g GuideCategory) IsValid() bool {
	switch g {
	case GuideCategoryFirstAid, GuideCategoryShelter, GuideCategoryFood, GuideCategoryWater, GuideCategoryNavigation,
		GuideCategoryFire, GuideCategorySignaling, GuideCategoryWeather, GuideCategoryWildlife, GuideCategoryTools:
		return true
	default:
		return false
	}
}

func (g *GuideCategory) Scan(value interface{}) error {
	if err := scanEnum(g, value); err != nil {
		return err
	}

	return validateEnum(*g, "GuideCategory", GuideCategory.IsValid)
}

func (g GuideCategory) Value() (driver.Value, error) {
	if err := validateEnum(g, "GuideCategory", GuideCategory.IsValid); err != nil {
		return nil, err
	}
	return string(g), nil
}

// ==================== GuideDifficulty ====================

// GuideDifficulty represents guide difficulty enum

type GuideDifficulty string

const (
	GuideDifficultyEasy   GuideDifficulty = "easy"
	GuideDifficultyMedium GuideDifficulty = "medium"
	GuideDifficultyHard   GuideDifficulty = "hard"
)

func (g GuideDifficulty) IsValid() bool {
	switch g {
	case GuideDifficultyEasy, GuideDifficultyMedium, GuideDifficultyHard:
		return true
	default:
		return false
	}
}

func (g *GuideDifficulty) Scan(value interface{}) error {
	if err := scanEnum(g, value); err != nil {
		return err
	}

	return validateEnum(*g, "GuideDifficulty", GuideDifficulty.IsValid)
}

func (g GuideDifficulty) Value() (driver.Value, error) {
	if err := validateEnum(g, "GuideDifficulty", GuideDifficulty.IsValid); err != nil {
		return nil, err
	}

	return string(g), nil
}

// ==================== NotificationType ====================

// NotificationType represents notification type enum

type NotificationType string

const (
	NotificationTypeWeatherAlert NotificationType = "weather_alert"
	NotificationTypeDangerZone   NotificationType = "danger_zone"
	NotificationTypeEmergency    NotificationType = "emergency"
	NotificationTypeGroupMessage NotificationType = "group_message"
	NotificationTypeSystem       NotificationType = "system"
	NotificationTypeAchivement   NotificationType = "achivement"
	NotificationTypeSubscription NotificationType = "subscription"
)

func (n NotificationType) IsValid() bool {
	switch n {
	case NotificationTypeWeatherAlert, NotificationTypeDangerZone, NotificationTypeEmergency, NotificationTypeGroupMessage,
		NotificationTypeSystem, NotificationTypeAchivement, NotificationTypeSubscription:
		return true
	default:
		return false
	}
}

func (n *NotificationType) Scan(value interface{}) error {
	if err := scanEnum(n, value); err != nil {
		return err
	}

	return validateEnum(*n, "NotificationType", NotificationType.IsValid)
}

func (n NotificationType) Value() (driver.Value, error) {
	if err := validateEnum(n, "NotificationType", NotificationType.IsValid); err != nil {
		return nil, err
	}

	return string(n), nil
}

// ==================== SubscriptionStatus ====================

// SubscriptionStatus represents subscription status enum
type SubscriptionStatus string

const (
	SubscriptionStatusTrial       SubscriptionStatus = "trial"
	SubscriptionStatusActive      SubscriptionStatus = "active"
	SubscriptionStatusExpired     SubscriptionStatus = "expired"
	SubscriptionStatusCancelled   SubscriptionStatus = "cancelled"
	SubscriptionStatusPaused      SubscriptionStatus = "paused"
	SubscriptionStatusGracePeriod SubscriptionStatus = "grace_period"
	SubscriptionStatusRefunded    SubscriptionStatus = "refunded"
	SubscriptionStatusRevoked     SubscriptionStatus = "revoked"
)

func (s SubscriptionStatus) IsValid() bool {
	switch s {
	case SubscriptionStatusTrial, SubscriptionStatusActive, SubscriptionStatusExpired, SubscriptionStatusCancelled,
		SubscriptionStatusPaused, SubscriptionStatusGracePeriod, SubscriptionStatusRefunded, SubscriptionStatusRevoked:
		return true
	default:
		return false
	}
}

func (s *SubscriptionStatus) Scan(value interface{}) error {
	if err := scanEnum(s, value); err != nil {
		return err
	}

	return validateEnum(*s, "SubscriptionStatus", SubscriptionStatus.IsValid)
}

func (s SubscriptionStatus) Value() (driver.Value, error) {
	if err := validateEnum(s, "SubscriptionStatus", SubscriptionStatus.IsValid); err != nil {
		return nil, err
	}

	return string(s), nil
}

// ==================== PlatformType ====================

// PlatformType represents platform type enum

type PlatformType string

const (
	PlatformTypeIOS     PlatformType = "ios"
	PlatformTypeAndroid PlatformType = "android"
	PlatformTypeWeb     PlatformType = "web"
)

func (p PlatformType) IsValid() bool {
	switch p {
	case PlatformTypeIOS, PlatformTypeAndroid, PlatformTypeWeb:
		return true
	default:
		return false

	}
}

func (p *PlatformType) Scan(value interface{}) error {
	if err := scanEnum(p, value); err != nil {
		return err
	}

	return validateEnum(*p, "PlatformType", PlatformType.IsValid)
}

func (p PlatformType) Value() (driver.Value, error) {
	if err := validateEnum(p, "PlatformType", PlatformType.IsValid); err != nil {
		return nil, err
	}

	return string(p), nil
}

// ==================== TransactionType ====================

// TransactionType represents transaction type enum

type TransactionType string

const (
	TransactionTypePurchase     TransactionType = "purchase"
	TransactionTypeRenewal      TransactionType = "renewal"
	TransactionTypeRefund       TransactionType = "refund"
	TransactionTypeUpgrade      TransactionType = "upgrade"
	TransactionTypeDowngrade    TransactionType = "downgrade"
	TransactionTypeCancellation TransactionType = "cancellation"
	TransactionTypeAdminGrant   TransactionType = "admin_grant"
	TransactionTypeAdminRevoke  TransactionType = "admin_revoke"
)

func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypePurchase, TransactionTypeRenewal, TransactionTypeRefund, TransactionTypeUpgrade,
		TransactionTypeDowngrade, TransactionTypeCancellation, TransactionTypeAdminGrant, TransactionTypeAdminRevoke:
		return true
	default:
		return false
	}
}

func (t *TransactionType) Scan(value interface{}) error {
	if err := scanEnum(t, value); err != nil {
		return err
	}
	return validateEnum(*t, "TransactionType", TransactionType.IsValid)
}

func (t TransactionType) Value() (driver.Value, error) {
	if err := validateEnum(t, "TransactionType", TransactionType.IsValid); err != nil {
		return nil, err
	}
	return string(t), nil
}

// ==================== TransactionStatus ====================

// TransactionStatus represents transaction status enum
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

func (t TransactionStatus) IsValid() bool {
	switch t {
	case TransactionStatusPending, TransactionStatusCompleted, TransactionStatusFailed, TransactionStatusRefunded:
		return true
	default:
		return false
	}
}

func (t *TransactionStatus) Scan(value interface{}) error {
	if err := scanEnum(t, value); err != nil {
		return err
	}
	return validateEnum(*t, "TransactionStatus", TransactionStatus.IsValid)
}

func (t TransactionStatus) Value() (driver.Value, error) {
	if err := validateEnum(t, "TransactionStatus", TransactionStatus.IsValid); err != nil {
		return nil, err
	}
	return string(t), nil
}

// ==================== DangerSeverity ====================

// DangerSeverity represents danger severity enum

type DangerSeverity string

const (
	DangerSeverityLow      DangerSeverity = "low"
	DangerSeverityMedium   DangerSeverity = "medium"
	DangerSeverityHigh     DangerSeverity = "high"
	DangerSeverityCritical DangerSeverity = "critical"
)

func (d DangerSeverity) IsValid() bool {
	switch d {
	case DangerSeverityLow, DangerSeverityMedium, DangerSeverityHigh, DangerSeverityCritical:
		return true
	default:
		return false
	}
}

func (d *DangerSeverity) Scan(value interface{}) error {
	if err := scanEnum(d, value); err != nil {
		return err
	}

	return validateEnum(*d, "DangerSeverity", DangerSeverity.IsValid)
}

func (d DangerSeverity) Value() (driver.Value, error) {
	if err := validateEnum(d, "DangerSeverity", DangerSeverity.IsValid); err != nil {
		return nil, err
	}

	return string(d), nil
}

// ==================== DangerType ====================

// DangerType represents danger type enum
type DangerType string

const (
	DangerTypeFlood      DangerType = "flood"
	DangerTypeStorm      DangerType = "storm"
	DangerTypeEarthquake DangerType = "earthquake"
	DangerTypeFire       DangerType = "fire"
	DangerTypeLandslide  DangerType = "landslide"
	DangerTypeVolcano    DangerType = "volcano"
	DangerTypeConflict   DangerType = "conflict"
	DangerTypeRadiation  DangerType = "radiation"
	DangerTypeOther      DangerType = "other"
)

func (d DangerType) IsValid() bool {
	switch d {
	case DangerTypeFlood, DangerTypeStorm, DangerTypeEarthquake, DangerTypeFire, DangerTypeLandslide,
		DangerTypeVolcano, DangerTypeConflict, DangerTypeRadiation, DangerTypeOther:
		return true
	default:
		return false
	}
}

func (d *DangerType) Scan(value interface{}) error {
	if err := scanEnum(d, value); err != nil {
		return err
	}

	return validateEnum(*d, "DangerType", DangerType.IsValid)
}

func (d DangerType) Value() (driver.Value, error) {
	if err := validateEnum(d, "DangerType", DangerType.IsValid); err != nil {
		return nil, err
	}
	return string(d), nil
}

// ==================== MemberStatus ====================

// MemberStatus represents group member status enum

type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusInactive MemberStatus = "inactive"
	MemberStatusRemoved  MemberStatus = "removed"
)

func (m MemberStatus) IsValid() bool {
	switch m {
	case MemberStatusActive, MemberStatusInactive, MemberStatusRemoved:
		return true
	default:
		return false
	}
}

func (m *MemberStatus) Scan(value interface{}) error {
	if err := scanEnum(m, value); err != nil {
		return err
	}

	return validateEnum(*m, "MemberStatus", MemberStatus.IsValid)
}

func (m MemberStatus) Value() (driver.Value, error) {
	if err := validateEnum(m, "MemberStatus", MemberStatus.IsValid); err != nil {
		return nil, err
	}

	return string(m), nil
}

// ==================== SOSStatus ====================

// SOSTrigger represents how SOS was triggered
type SOSTrigger string

const (
	SOSTriggerManual        SOSTrigger = "manual"
	SOSTriggerAutoCrash     SOSTrigger = "auto_crash"
	SOSTriggerFallDetection SOSTrigger = "fall_detection"
	SOSTriggerPanicButton   SOSTrigger = "panic_button"
)

func (s SOSTrigger) IsValid() bool {
	switch s {
	case SOSTriggerManual, SOSTriggerAutoCrash, SOSTriggerFallDetection, SOSTriggerPanicButton:
		return true
	default:
		return false
	}
}

// SOSStatus represents SOS event status
type SOSStatus string

const (
	SOSStatusActive     SOSStatus = "active"
	SOSStatusResponding SOSStatus = "responding"
	SOSStatusResolved   SOSStatus = "resolved"
	SOSStatusFalseAlarm SOSStatus = "false_alarm"
	SOSStatusCancelled  SOSStatus = "cancelled"
)

func (s SOSStatus) IsValid() bool {
	switch s {
	case SOSStatusActive, SOSStatusResponding, SOSStatusResolved,
		SOSStatusFalseAlarm, SOSStatusCancelled:
		return true
	default:
		return false
	}
}

func (s *SOSStatus) Scan(value interface{}) error {
	if err := scanEnum(s, value); err != nil {
		return err
	}

	return validateEnum(*s, "SOSStatus", SOSStatus.IsValid)
}

func (s SOSStatus) Value() (driver.Value, error) {
	if err := validateEnum(s, "SOSStatus", SOSStatus.IsValid); err != nil {
		return nil, err
	}

	return string(s), nil
}

// ==================== Point (Geographic) ====================

// Point represents a geographic point (longitude, latitude)

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// IsValid checks if the point has valid coordinates
func (p *Point) IsValid() bool {
	if p == nil {
		return false
	}
	return p.Latitude >= -90 && p.Latitude <= 90 &&
		p.Longitude >= -180 && p.Longitude <= 180
}

func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot scan %T into Point", value)
	}

	// Parse POINT string: "POINT(longitude latitude)"
	var lon, lat float64
	_, err := fmt.Sscan(str, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return fmt.Errorf("invalid POINT format: %w", err)
	}

	p.Longitude = lon
	p.Latitude = lat
	return nil
}

// Value implements driver.Valuer for PostGIS geography point
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%f %f)", p.Longitude, p.Latitude), nil
}

// ==================== NullableString ====================

// NullableString is a helper for nullable strings

type NullableString struct {
	String string
	Valid  bool
}

func (ns *NullableString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}

	switch v := value.(type) {
	case string:
		ns.String = v
	case []byte:
		ns.String = string(v)
	default:
		return fmt.Errorf("cannot scan %T into NullableString", value)
	}

	ns.Valid = true
	return nil
}

func (ns NullableString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}

	return ns.String, nil
}

// ==================== Custom Types ====================

// JSONB is a custom type for PostgreSQL JSONB
type JSONB map[string]any

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONB: %T", value)
	}
	return json.Unmarshal(data, j)
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil || len(j) == 0 {
		return []byte(`{}`), nil
	}
	return json.Marshal(j)
}

// GeographyPolygon represents PostGIS geography polygon raw data
type GeographyPolygon []byte

// Scan implements sql.Scanner for PostGIS geography
func (g *GeographyPolygon) Scan(value interface{}) error {
	if value == nil {
		*g = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*g = v
	case string:
		*g = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into GeographyPolygon", value)
	}
	return nil
}

// Value implements driver.Valuer for PostGIS
func (g GeographyPolygon) Value() (driver.Value, error) {
	if g == nil {
		return nil, nil
	}
	return []byte(g), nil
}

// ToPolygon converts raw PostGIS data to Polygon struct
func (g GeographyPolygon) ToPolygon() (*Polygon, error) {
	if g == nil || len(g) == 0 {
		return nil, nil
	}

	var polygon Polygon
	if err := json.Unmarshal(g, &polygon); err != nil {
		return nil, fmt.Errorf("failed to parse geography to polygon: %w", err)
	}

	return &polygon, nil
}

// Polygon represents a geographic polygon for danger zones (GeoJSON format)
type Polygon struct {
	Type        string        `json:"type"`        // Always "Polygon"
	Coordinates [][][]float64 `json:"coordinates"` // [[[lon,lat], [lon,lat], ...]]
}

// IsValid checks if polygon has valid GeoJSON structure
func (p *Polygon) IsValid() bool {
	if p == nil || p.Type != "Polygon" {
		return false
	}

	if len(p.Coordinates) == 0 {
		return false
	}

	// Polygon must have at least one ring with at least 4 points
	// First and last point must be the same (closed ring)
	for _, ring := range p.Coordinates {
		if len(ring) < 4 {
			return false
		}
		// Check if first and last points are the same
		first := ring[0]
		last := ring[len(ring)-1]
		if len(first) != 2 || len(last) != 2 {
			return false
		}
		if first[0] != last[0] || first[1] != last[1] {
			return false
		}
		// Validate longitude and latitude ranges
		for _, coord := range ring {
			if len(coord) != 2 {
				return false
			}
			lon, lat := coord[0], coord[1]
			if lon < -180 || lon > 180 || lat < -90 || lat > 90 {
				return false
			}
		}
	}
	return true
}

// ToGeoJSON converts Polygon to GeoJSON string for PostGIS
func (p *Polygon) ToGeoJSON() (string, error) {
	if p == nil {
		return "", fmt.Errorf("polygon is nil")
	}

	geoJSON, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("failed to marshal polygon: %w", err)
	}

	return string(geoJSON), nil
}

// SafetyLevel represents safety level enum
type SafetyLevel string

const (
	SafetyLevelSafe     SafetyLevel = "safe"
	SafetyLevelWarning  SafetyLevel = "warning"
	SafetyLevelDanger   SafetyLevel = "danger"
	SafetyLevelCritical SafetyLevel = "critical"
)

func (s SafetyLevel) IsValid() bool {
	switch s {
	case SafetyLevelSafe, SafetyLevelWarning, SafetyLevelDanger, SafetyLevelCritical:
		return true
	default:
		return false
	}
}

// AdminAction represents admin action type enum
type AdminAction string

const (
	// User actions
	AdminActionUserActivate   AdminAction = "user_activate"
	AdminActionUserDeactivate AdminAction = "user_deactivate"
	AdminActionUserPromote    AdminAction = "user_promote"
	AdminActionUserDemote     AdminAction = "user_demote"
	AdminActionUserDelete     AdminAction = "user_delete"
	AdminActionUserBan        AdminAction = "user_ban"
	AdminActionUserUnban      AdminAction = "user_unban"

	// Group actions
	AdminActionGroupActivate       AdminAction = "group_activate"
	AdminActionGroupDeactivate     AdminAction = "group_deactivate"
	AdminActionGroupDelete         AdminAction = "group_delete"
	AdminActionGroupTransferLeader AdminAction = "group_transfer_leader"

	// SOS actions
	AdminActionSOSResolve    AdminAction = "sos_resolve"
	AdminActionSOSFalseAlarm AdminAction = "sos_false_alarm"

	// Danger zone actions
	AdminActionDangerZoneCreate AdminAction = "danger_zone_create"
	AdminActionDangerZoneUpdate AdminAction = "danger_zone_update"
	AdminActionDangerZoneDelete AdminAction = "danger_zone_delete"

	// Subscription actions
	AdminActionSubscriptionGrant  AdminAction = "subscription_grant"
	AdminActionSubscriptionRevoke AdminAction = "subscription_revoke"
	AdminActionSubscriptionExtend AdminAction = "subscription_extend"

	// System actions
	AdminActionSystemConfig      AdminAction = "system_config"
	AdminActionSystemMaintenance AdminAction = "system_maintenance"
)

func (a AdminAction) IsValid() bool {
	switch a {
	case AdminActionUserActivate, AdminActionUserDeactivate, AdminActionUserPromote, AdminActionUserDemote,
		AdminActionUserDelete, AdminActionUserBan, AdminActionUserUnban,
		AdminActionGroupActivate, AdminActionGroupDeactivate, AdminActionGroupDelete, AdminActionGroupTransferLeader,
		AdminActionSOSResolve, AdminActionSOSFalseAlarm,
		AdminActionDangerZoneCreate, AdminActionDangerZoneUpdate, AdminActionDangerZoneDelete,
		AdminActionSubscriptionGrant, AdminActionSubscriptionRevoke, AdminActionSubscriptionExtend,
		AdminActionSystemConfig, AdminActionSystemMaintenance:
		return true
	default:
		return false
	}
}

// AuditTargetType represents the type of target in audit logs
type AuditTargetType string

const (
	AuditTargetUser         AuditTargetType = "user"
	AuditTargetGroup        AuditTargetType = "group"
	AuditTargetSOSEvent     AuditTargetType = "sos_event"
	AuditTargetDangerZone   AuditTargetType = "danger_zone"
	AuditTargetSubscription AuditTargetType = "subscription"
	AuditTargetSystem       AuditTargetType = "system"
)

func (a AuditTargetType) IsValid() bool {
	switch a {
	case AuditTargetUser, AuditTargetGroup, AuditTargetSOSEvent,
		AuditTargetDangerZone, AuditTargetSubscription, AuditTargetSystem:
		return true
	default:
		return false
	}
}

// HealthStatus represents system health status enum
type HealthStatus string

const (
	HealthStatusHealthy  HealthStatus = "healthy"
	HealthStatusDegraded HealthStatus = "degraded"
	HealthStatusDown     HealthStatus = "down"
)

func (h HealthStatus) IsValid() bool {
	switch h {
	case HealthStatusHealthy, HealthStatusDegraded, HealthStatusDown:
		return true
	default:
		return false
	}
}

// AdminAuditLog represents an admin action audit log
type AdminAuditLog struct {
	ID           int              `json:"id" db:"id"`
	AdminID      *int             `json:"admin_id,omitempty" db:"admin_id"`
	Action       AdminAction      `json:"action" db:"action"`
	TargetType   *AuditTargetType `json:"target_type,omitempty" db:"target_type"`
	TargetID     *int             `json:"target_id,omitempty" db:"target_id"`
	Reason       *string          `json:"reason,omitempty" db:"reason"`
	Metadata     *JSONB           `json:"metadata,omitempty" db:"metadata"` // Pointer để nullable và omitempty hoạt động
	Success      bool             `json:"success" db:"success"`
	ErrorMessage *string          `json:"error_message,omitempty" db:"error_message"`
	IPAddress    *string          `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string          `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
}

// AdminAuditLogCreate represents audit log creation request
type AdminAuditLogCreate struct {
	Action       AdminAction      `json:"action" binding:"required"`
	TargetType   *AuditTargetType `json:"target_type,omitempty"`
	TargetID     *int             `json:"target_id,omitempty" binding:"omitempty,min=1"`
	Reason       *string          `json:"reason,omitempty" binding:"omitempty,max=1000"`
	Metadata     map[string]any   `json:"metadata,omitempty"` // Go 1.18+
	Success      bool             `json:"success"`
	ErrorMessage *string          `json:"error_message,omitempty" binding:"omitempty,max=2000"`
	IPAddress    *string          `json:"ip_address,omitempty" binding:"omitempty,ip"`
	UserAgent    *string          `json:"user_agent,omitempty" binding:"omitempty,max=500"`
}

// Validate validates the audit log create request
func (a *AdminAuditLogCreate) Validate() error {
	if !a.Action.IsValid() {
		return fmt.Errorf("invalid admin action")
	}

	if a.TargetType != nil && !a.TargetType.IsValid() {
		return fmt.Errorf("invalid target type")
	}

	return nil
}
