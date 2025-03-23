package main

import (
	"car-service/cmd/api/controllers"
	api "car-service/cmd/api/mediator"
	"car-service/cmd/api/server"
	"car-service/internal/application/commands"
	"car-service/internal/application/queries"
	"car-service/internal/application/services"
	"car-service/internal/domain/repositories"
	gormrepo "car-service/internal/infrastructure/gorm"
	"car-service/internal/infrastructure/migrations"
	"car-service/pkg/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Cargar configuración de ambiente
	env, err := config.NewEnvironment()
	if err != nil {
		return err
	}

	// Configuración de la base de datos
	db, err := setupDatabase(env)
	if err != nil {
		return err
	}

	// Inicializar repositorios
	var carRepo repositories.CarRepository = gormrepo.NewCarRepository(db)
	carService := services.NewCarService(carRepo)

	mediator := api.NewMediator()

	newCarCommand := commands.NewNewCarCommand(carService)
	mediator.Register("CreateCar", newCarCommand)
	mediator.RegisterQuery("GetCars", queries.NewGetCarsQuery(carService))

	carController := controllers.NewCarController(mediator)

	// Configurar el servidor
	serverCfg := &server.ServerConfig{
		CarController: carController,
		Port:          env.ServerPort,
	}

	// Crear y configurar el servidor
	srv := server.NewServer(serverCfg)

	// Iniciar el servidor (bloqueante)
	return srv.Start()
}
func setupDatabase(env *config.Environment) (*gorm.DB, error) {
	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(env.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ejecutar migraciones
	if err := migrations.Migrate(db); err != nil {
		return nil, err
	}
	log.Println("Migraciones ejecutadas correctamente")

	return db, nil
}
