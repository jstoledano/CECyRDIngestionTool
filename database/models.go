package database

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Folio                             string `gorm:"primaryKey;index;unique"`
	Estatus                           string
	CausaRechazo                      string
	MovimientoSolicitado              string
	MovimientoDefinitivo              string
	FechaTramite                      time.Time
	FechaRecibidoCecyrd               time.Time
	FechaRegistradoCecyrd             time.Time
	FechaRechazado                    time.Time
	FechaCanceladoMovimientoPosterior time.Time
	FechaAltaPe                       time.Time
	FechaAfectacionPadron             time.Time
	FechaActualizacionPe              time.Time
	FechaReincorporacionPe            time.Time
	FechaExitoso                      time.Time
	FechaLoteProduccion               time.Time
	FechaListoReimpresion             time.Time
	FechaCpvCreada                    time.Time
	FechaCpvRegistradaMac             time.Time
	FechaCpvDisponible                time.Time
	FechaCpvEntregada                 time.Time
	FechaAfectacionLn                 time.Time
	Distrito                          int    `gorm:"not null"`
	Mac                               string `gorm:"type:varchar(6);collate:pg_catalog.default;not null"`
	TramoDisponible                   time.Duration
	TramoEntrega                      time.Duration
	TramoExitoso                      time.Duration
}
