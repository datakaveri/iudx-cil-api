package models

import (
	"iudx_domain_specific_apis/pkg/db"
	"time"
)

type AQMSpatialForecast struct {
	GeoJson             string    `db:"geojson" json:"geojson"`
	ObservationDateTime time.Time `db:"observationdatetime" json:"observationdatetime"`
	PollutantVal        float32   `db:"pollutant_val" json:"pollutant_val"`
}

type AQMSpatialForecastMinMax struct {
	Min float32 `db:"min" json:"min"`
	Max float32 `db:"max" json:"max"`
}

type AQMSpatialForecastModel struct{}

func (m AQMSpatialForecastModel) GetAQMSpatialForecasts(startTime time.Time, endTime time.Time, measuredValue string) (aqmSpatialForecast []AQMSpatialForecast, aqmSpatialForecastMinMax []AQMSpatialForecastMinMax, err error) {

	_, err = db.GetDB().Select(&aqmSpatialForecast, `
		select kdmc_aqm_geojsons.geojson,
		kdmc_aqm_interpolation_actual_data.observationdatetime, kdmc_aqm_interpolation_actual_data.`+measuredValue+` as pollutant_val
		from kdmc_aqm_geojsons
		inner join kdmc_aqm_interpolation_actual_data on kdmc_aqm_geojsons.hex_id = kdmc_aqm_interpolation_actual_data.hex_id
		where observationdatetime >= $1 and observationdatetime <= $2
		order by kdmc_aqm_interpolation_actual_data.observationdatetime
	`, startTime, endTime)

	_, err = db.GetDB().Select(&aqmSpatialForecastMinMax, `
		select min(kdmc_aqm_interpolation_actual_data.`+measuredValue+`), 
		max(kdmc_aqm_interpolation_actual_data.`+measuredValue+`) 
		from kdmc_aqm_interpolation_actual_data
		where observationdatetime >= $1 and observationdatetime <= $2
	`, startTime, endTime)

	return aqmSpatialForecast, aqmSpatialForecastMinMax, err
}
