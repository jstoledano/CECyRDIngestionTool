package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jstoledano/CECyRDIngestionTool/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Open the csv file
	f, err := os.Open("data/01.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(f)
	csvReader.Comma = '|'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	loc, _ := time.LoadLocation("America/Mexico_City")
	layout := "02/01/2006 03:04:05 PM"

	for _, row := range data {
		t := database.Record{}
		t.Folio = row[0]
		t.Estatus = row[1]
		t.CausaRechazo = row[2]
		t.MovimientoSolicitado = row[3]
		t.MovimientoDefinitivo = row[4]
		t.FechaTramite, _ = time.ParseInLocation(layout, row[5], loc)
		t.FechaRecibidoCecyrd, _ = time.ParseInLocation(layout, row[6], loc)
		t.FechaRegistradoCecyrd, _ = time.ParseInLocation(layout, row[7], loc)
		t.FechaRechazado, _ = time.ParseInLocation(layout, row[8], loc)
		t.FechaCanceladoMovimientoPosterior, _ = time.ParseInLocation(layout, row[9], loc)
		t.FechaAltaPe, _ = time.ParseInLocation(layout, row[10], loc)
		t.FechaAfectacionPadron, _ = time.ParseInLocation(layout, row[11], loc)
		t.FechaActualizacionPe, _ = time.ParseInLocation(layout, row[12], loc)
		t.FechaReincorporacionPe, _ = time.ParseInLocation(layout, row[13], loc)
		t.FechaExitoso, _ = time.ParseInLocation(layout, row[14], loc)
		t.FechaLoteProduccion, _ = time.ParseInLocation(layout, row[15], loc)
		t.FechaListoReimpresion, _ = time.ParseInLocation(layout, row[16], loc)
		t.FechaCpvCreada, _ = time.ParseInLocation(layout, row[17], loc)
		t.FechaCpvRegistradaMac, _ = time.ParseInLocation(layout, row[18], loc)
		t.FechaCpvDisponible, _ = time.ParseInLocation(layout, row[19], loc)
		t.FechaCpvEntregada, _ = time.ParseInLocation(layout, row[20], loc)
		t.FechaAfectacionLn, _ = time.ParseInLocation(layout, row[21], loc)
		t.Distrito, _ = strconv.Atoi(t.Folio[5:6])
		t.Mac = t.Folio[2:8]
		t.TramoDisponible = t.FechaCpvDisponible.Sub(t.FechaTramite)
		t.TramoEntrega = t.FechaCpvEntregada.Sub(t.FechaCpvDisponible)
		t.TramoExitoso = t.FechaExitoso.Sub(t.FechaTramite)

		err := upsertRecord(db, t)
		if err != nil {
			log.Fatal("Error al procesar ", t.Folio, " -- ", err)
		}
	}
	defer f.Close()
}

func upsertRecord(db *gorm.DB, t database.Record) error {
	log.Println("Processing record", t.Folio)
	tx := db.Begin()
	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "folio"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"estatus",
			"causa_rechazo",
			"movimiento_solicitado",
			"movimiento_definitivo",
			"fecha_tramite",
			"fecha_recibido_cecyrd",
			"fecha_registrado_cecyrd",
			"fecha_rechazado",
			"fecha_cancelado_movimiento_posterior",
			"fecha_alta_pe",
			"fecha_afectacion_padron",
			"fecha_actualizacion_pe",
			"fecha_reincorporacion_pe",
			"fecha_exitoso",
			"fecha_lote_produccion",
			"fecha_listo_reimpresion",
			"fecha_cpv_creada",
			"fecha_cpv_registrada_mac",
			"fecha_cpv_disponible",
			"fecha_cpv_entregada",
			"fecha_afectacion_ln",
			"distrito",
			"mac",
			"tramo_disponible",
			"tramo_entrega",
			"tramo_exitoso",
		}),
	}).Create(&t).Error

	if err != nil {
		tx.Rollback()
		log.Println("Error al procesar ", t.Folio, " -- ", err)
		return err
	}

	log.Println("Record", t.Folio, "processed")
	tx.Commit()
	return nil
}
