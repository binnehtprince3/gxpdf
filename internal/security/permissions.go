package security

// Permission represents PDF document permissions.
//
// These flags control what operations are allowed on an encrypted PDF.
// Multiple permissions can be combined using the OR operator (|).
//
// Example:
//
//	perms := PermissionPrint | PermissionCopy | PermissionModify
type Permission int32

const (
	// PermissionPrint allows printing the document (bit 3).
	PermissionPrint Permission = 1 << 2

	// PermissionModify allows modifying the document (bit 4).
	PermissionModify Permission = 1 << 3

	// PermissionCopy allows copying text and graphics (bit 5).
	PermissionCopy Permission = 1 << 4

	// PermissionAnnotate allows adding or modifying annotations (bit 6).
	PermissionAnnotate Permission = 1 << 5

	// PermissionFillForms allows filling form fields (bit 9).
	PermissionFillForms Permission = 1 << 8

	// PermissionExtract allows extracting text for accessibility (bit 10).
	PermissionExtract Permission = 1 << 9

	// PermissionAssemble allows assembling the document (bit 11).
	PermissionAssemble Permission = 1 << 10

	// PermissionPrintHighQuality allows high-quality printing (bit 12).
	PermissionPrintHighQuality Permission = 1 << 11

	// PermissionAll grants all permissions.
	PermissionAll Permission = PermissionPrint |
		PermissionModify |
		PermissionCopy |
		PermissionAnnotate |
		PermissionFillForms |
		PermissionExtract |
		PermissionAssemble |
		PermissionPrintHighQuality

	// PermissionNone grants no permissions (default).
	PermissionNone Permission = 0
)

// Has checks if a specific permission is granted.
//
// Example:
//
//	perms := PermissionPrint | PermissionCopy
//	if perms.Has(PermissionPrint) {
//	    fmt.Println("Printing allowed")
//	}
func (p Permission) Has(perm Permission) bool {
	return p&perm == perm
}

// Add adds a permission to the current permissions.
//
// Example:
//
//	perms := PermissionPrint
//	perms = perms.Add(PermissionCopy)
func (p Permission) Add(perm Permission) Permission {
	return p | perm
}

// Remove removes a permission from the current permissions.
//
// Example:
//
//	perms := PermissionAll
//	perms = perms.Remove(PermissionModify)
func (p Permission) Remove(perm Permission) Permission {
	return p &^ perm
}

// ToPDFValue converts permissions to the PDF integer format.
//
// The PDF specification requires bits 1, 2, 7, 8 to be set to 1,
// and all other bits depend on the actual permissions.
func (p Permission) ToPDFValue() int32 {
	// Set required bits (1, 2, 7, 8) to 1.
	// Bits are 0-indexed: bit 1 = 0x02, bit 2 = 0x04, bit 7 = 0x80, bit 8 = 0x100.
	const requiredBits int32 = 0x02 | 0x04 | 0x80 | 0x100

	// Combine with actual permissions.
	// Use bitwise OR and then set all high bits to 1 (PDF spec requirement).
	result := int32(p) | requiredBits

	// Set all bits above bit 12 to 1 (PDF spec: complement low bits).
	// This is done by OR-ing with 0xFFFFF000 (all bits above 12 are 1).
	result |= ^int32(0xFFF)

	return result
}

// String returns a human-readable string of enabled permissions.
func (p Permission) String() string {
	if p == PermissionNone {
		return "None"
	}

	if p == PermissionAll {
		return "All"
	}

	return buildPermissionString(p)
}

// buildPermissionString builds the permission string from individual flags.
func buildPermissionString(p Permission) string {
	perms := collectPermissions(p)

	result := ""
	for i, perm := range perms {
		if i > 0 {
			result += " | "
		}
		result += perm
	}

	return result
}

// collectPermissions collects all enabled permission names.
func collectPermissions(p Permission) []string {
	var perms []string

	permChecks := []struct {
		flag Permission
		name string
	}{
		{PermissionPrint, "Print"},
		{PermissionModify, "Modify"},
		{PermissionCopy, "Copy"},
		{PermissionAnnotate, "Annotate"},
		{PermissionFillForms, "FillForms"},
		{PermissionExtract, "Extract"},
		{PermissionAssemble, "Assemble"},
		{PermissionPrintHighQuality, "PrintHighQuality"},
	}

	for _, pc := range permChecks {
		if p.Has(pc.flag) {
			perms = append(perms, pc.name)
		}
	}

	return perms
}
