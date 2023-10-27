package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
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
	f, err := os.Open("data/01.Trámites_Tlax_01Septiembre-30Noviembre2021.txt")
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

	log.Println("Processing", len(data), "records")
	bar := pb.StartNew(len(data))
	start := time.Now()
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
		bar.Increment()
		if err != nil {
			log.Fatal("Error al procesar ", t.Folio, " -- ", err)
		}
	}
	defer f.Close()

	bar.Finish()
	end := time.Now()

	log.Println("Time elapsed:", end.Sub(start))
}

func upsertRecord(db *gorm.DB, t database.Record) error {
	tx := db.Begin()

	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "folio"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"estatus":                              t.Estatus,
			"causa_rechazo":                        t.CausaRechazo,
			"movimiento_solicitado":                t.MovimientoSolicitado,
			"movimiento_definitivo":                t.MovimientoDefinitivo,
			"fecha_tramite":                        t.FechaTramite,
			"fecha_recibido_cecyrd":                t.FechaRecibidoCecyrd,
			"fecha_registrado_cecyrd":              t.FechaRegistradoCecyrd,
			"fecha_rechazado":                      t.FechaRechazado,
			"fecha_cancelado_movimiento_posterior": t.FechaCanceladoMovimientoPosterior,
			"fecha_alta_pe":                        t.FechaAltaPe,
			"fecha_afectacion_padron":              t.FechaAfectacionPadron,
			"fecha_actualizacion_pe":               t.FechaActualizacionPe,
			"fecha_reincorporacion_pe":             t.FechaReincorporacionPe,
			"fecha_exitoso":                        t.FechaExitoso,
			"fecha_lote_produccion":                t.FechaLoteProduccion,
			"fecha_listo_reimpresion":              t.FechaListoReimpresion,
			"fecha_cpv_creada":                     t.FechaCpvCreada,
			"fecha_cpv_registrada_mac":             t.FechaCpvRegistradaMac,
			"fecha_cpv_disponible":                 t.FechaCpvDisponible,
			"fecha_cpv_entregada":                  t.FechaCpvEntregada,
			"fecha_afectacion_ln":                  t.FechaAfectacionLn,
			"distrito":                             t.Distrito,
			"mac":                                  t.Mac,
			"tramo_disponible":                     t.TramoDisponible,
			"tramo_entrega":                        t.TramoEntrega,
			"tramo_exitoso":                        t.TramoExitoso}),
	}).Create(&t).Error; err != nil {
		tx.Rollback()
		log.Println("Error al procesar ", t.Folio, " -- ", err)
	}
	tx.Commit()
	return nil
}
