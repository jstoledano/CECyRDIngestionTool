package database

import (
	"time"
)

type Tramite struct {
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
	Distrito                          int
	Mac                               string
	TramoDisponible                   time.Duration
	TramoEntrega                      time.Duration
	TramoExitoso                      time.Duration
}
