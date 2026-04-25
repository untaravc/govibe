package seeder

import "testing"

func TestLoadRegionAssets(t *testing.T) {
	provinces, err := loadProvinces("../../resources/assets/provinsi/provinsi.json")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(provinces); got == 0 {
		t.Fatal("expected provinces")
	}
	if provinces[0].ID != 11 || provinces[0].Code != "11" || provinces[0].Name != "ACEH" {
		t.Fatalf("unexpected first province: %+v", provinces[0])
	}

	cities, err := loadCities("../../resources/assets/kabupaten_kota")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(cities); got == 0 {
		t.Fatal("expected cities")
	}
	if cities[0].ID != 1101 || cities[0].ProvinceID != 11 || cities[0].Code != "1101" {
		t.Fatalf("unexpected first city: %+v", cities[0])
	}

	districts, err := loadDistricts("../../resources/assets/kecamatan")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(districts); got == 0 {
		t.Fatal("expected districts")
	}
	if districts[0].ID != 1101010 || districts[0].CityID != 1101 || districts[0].Code != "1101010" {
		t.Fatalf("unexpected first district: %+v", districts[0])
	}

	villages, err := loadVillages("../../resources/assets/kelurahan_desa")
	if err != nil {
		t.Fatal(err)
	}
	if got := len(villages); got == 0 {
		t.Fatal("expected villages")
	}
	if villages[0].ID != 1101010001 || villages[0].DistrictID != 1101010 || villages[0].Code != "1101010001" {
		t.Fatalf("unexpected first village: %+v", villages[0])
	}
}
