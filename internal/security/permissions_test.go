package security

import "testing"

func TestPermission_Has(t *testing.T) {
	tests := []struct {
		name  string
		perms Permission
		check Permission
		want  bool
	}{
		{
			name:  "has print permission",
			perms: PermissionPrint,
			check: PermissionPrint,
			want:  true,
		},
		{
			name:  "does not have modify permission",
			perms: PermissionPrint,
			check: PermissionModify,
			want:  false,
		},
		{
			name:  "has multiple permissions",
			perms: PermissionPrint | PermissionCopy,
			check: PermissionPrint,
			want:  true,
		},
		{
			name:  "has all permissions",
			perms: PermissionAll,
			check: PermissionPrint,
			want:  true,
		},
		{
			name:  "none has no permissions",
			perms: PermissionNone,
			check: PermissionPrint,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.perms.Has(tt.check); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermission_Add(t *testing.T) {
	tests := []struct {
		name  string
		perms Permission
		add   Permission
		want  Permission
	}{
		{
			name:  "add to none",
			perms: PermissionNone,
			add:   PermissionPrint,
			want:  PermissionPrint,
		},
		{
			name:  "add to existing",
			perms: PermissionPrint,
			add:   PermissionCopy,
			want:  PermissionPrint | PermissionCopy,
		},
		{
			name:  "add duplicate",
			perms: PermissionPrint,
			add:   PermissionPrint,
			want:  PermissionPrint,
		},
		{
			name:  "add multiple",
			perms: PermissionNone,
			add:   PermissionPrint | PermissionCopy,
			want:  PermissionPrint | PermissionCopy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.perms.Add(tt.add)
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermission_Remove(t *testing.T) {
	tests := []struct {
		name   string
		perms  Permission
		remove Permission
		want   Permission
	}{
		{
			name:   "remove from all",
			perms:  PermissionAll,
			remove: PermissionModify,
			want:   PermissionAll &^ PermissionModify,
		},
		{
			name:   "remove from single",
			perms:  PermissionPrint,
			remove: PermissionPrint,
			want:   PermissionNone,
		},
		{
			name:   "remove non-existent",
			perms:  PermissionPrint,
			remove: PermissionCopy,
			want:   PermissionPrint,
		},
		{
			name:   "remove from multiple",
			perms:  PermissionPrint | PermissionCopy,
			remove: PermissionPrint,
			want:   PermissionCopy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.perms.Remove(tt.remove)
			if got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermission_ToPDFValue(t *testing.T) {
	tests := []struct {
		name  string
		perms Permission
	}{
		{
			name:  "none permissions",
			perms: PermissionNone,
		},
		{
			name:  "print permission",
			perms: PermissionPrint,
		},
		{
			name:  "all permissions",
			perms: PermissionAll,
		},
		{
			name:  "mixed permissions",
			perms: PermissionPrint | PermissionCopy | PermissionModify,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfValue := tt.perms.ToPDFValue()

			// Verify required bits are set.
			// Bits 1, 2, 7, 8 must be 1.
			const requiredBits int32 = 0x02 | 0x04 | 0x80 | 0x100

			if pdfValue&requiredBits != requiredBits {
				t.Errorf("ToPDFValue() missing required bits: got %#x", pdfValue)
			}

			// Verify all bits above bit 12 are 1 (PDF spec requirement).
			const highBitsMask int32 = ^int32(0xFFF)
			if pdfValue&highBitsMask != highBitsMask {
				t.Errorf("ToPDFValue() high bits not set: got %#x", pdfValue)
			}

			// Verify permission bits are preserved.
			permBits := int32(tt.perms)
			if pdfValue&permBits != permBits {
				t.Errorf("ToPDFValue() lost permission bits: got %#x, want %#x", pdfValue, permBits)
			}
		})
	}
}

func TestPermission_String(t *testing.T) {
	tests := []struct {
		name  string
		perms Permission
		want  string
	}{
		{
			name:  "none",
			perms: PermissionNone,
			want:  "None",
		},
		{
			name:  "all",
			perms: PermissionAll,
			want:  "All",
		},
		{
			name:  "single permission",
			perms: PermissionPrint,
			want:  "Print",
		},
		{
			name:  "multiple permissions",
			perms: PermissionPrint | PermissionCopy,
			want:  "Print | Copy",
		},
		{
			name:  "all individual permissions",
			perms: PermissionPrint | PermissionModify | PermissionCopy | PermissionAnnotate | PermissionFillForms | PermissionExtract | PermissionAssemble | PermissionPrintHighQuality,
			want:  "All",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.perms.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionConstants(t *testing.T) {
	// Verify permission bit positions.
	tests := []struct {
		name     string
		perm     Permission
		wantBit  int
		wantFlag int32
	}{
		{"Print", PermissionPrint, 2, 1 << 2},
		{"Modify", PermissionModify, 3, 1 << 3},
		{"Copy", PermissionCopy, 4, 1 << 4},
		{"Annotate", PermissionAnnotate, 5, 1 << 5},
		{"FillForms", PermissionFillForms, 8, 1 << 8},
		{"Extract", PermissionExtract, 9, 1 << 9},
		{"Assemble", PermissionAssemble, 10, 1 << 10},
		{"PrintHighQuality", PermissionPrintHighQuality, 11, 1 << 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int32(tt.perm) != tt.wantFlag {
				t.Errorf("Permission %s = %#x, want %#x", tt.name, tt.perm, tt.wantFlag)
			}
		})
	}
}

func TestPermission_ChainedOperations(t *testing.T) {
	// Test chaining Add and Remove operations.
	perms := PermissionNone
	perms = perms.Add(PermissionPrint)
	perms = perms.Add(PermissionCopy)
	perms = perms.Add(PermissionModify)

	if !perms.Has(PermissionPrint) {
		t.Error("Should have Print permission")
	}
	if !perms.Has(PermissionCopy) {
		t.Error("Should have Copy permission")
	}
	if !perms.Has(PermissionModify) {
		t.Error("Should have Modify permission")
	}

	perms = perms.Remove(PermissionModify)

	if !perms.Has(PermissionPrint) {
		t.Error("Should still have Print permission")
	}
	if !perms.Has(PermissionCopy) {
		t.Error("Should still have Copy permission")
	}
	if perms.Has(PermissionModify) {
		t.Error("Should not have Modify permission")
	}
}
