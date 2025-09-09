package cmd

import (
	"fmt"
	"github.com/gloowl/simple_crud/src/internal/database"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	dbConfig database.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "herbs-cli",
	Short: "CLI для управления базой данных лекарственных трав",
	Long: `Приложение командной строки для выполнения CRUD операций 
с базой данных лекарственных трав.

Поддерживаемые операции:
- Создание новых записей о травах
- Просмотр информации о травах
- Обновление данных о травах
- Удаление записей о травах
- Поиск трав по названию`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Connect to database before running any command
		if err := database.Connect(dbConfig); err != nil {
			fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
			os.Exit(1)
		}
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Close database connection after running command
		if err := database.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с БД: %v", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "файл конфигурации (по умолчанию $HOME/.herbs-cli.yaml)")

	// Database connection flags (используем вашу конфигурацию по умолчанию)
	rootCmd.PersistentFlags().StringVar(&dbConfig.Host, "host", "localhost", "адрес сервера PostgreSQL")
	rootCmd.PersistentFlags().IntVar(&dbConfig.Port, "port", 5432, "порт PostgreSQL")
	rootCmd.PersistentFlags().StringVar(&dbConfig.User, "user", "admin", "имя пользователя PostgreSQL")
	rootCmd.PersistentFlags().StringVar(&dbConfig.Password, "password", "pwd4adm", "пароль PostgreSQL")
	rootCmd.PersistentFlags().StringVar(&dbConfig.DBName, "dbname", "simple_crud_db", "имя базы данных")
	rootCmd.PersistentFlags().StringVar(&dbConfig.SSLMode, "sslmode", "disable", "режим SSL (disable, require, verify-ca, verify-full)")

	// Bind flags to viper
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("dbname", rootCmd.PersistentFlags().Lookup("dbname"))
	viper.BindPFlag("sslmode", rootCmd.PersistentFlags().Lookup("sslmode"))
}

// initConfig reads in config file and ENV variables
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".herbs-cli" (without extension)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("herbs-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Используется конфигурационный файл: %s\n", viper.ConfigFileUsed())

		// Update database config from viper
		dbConfig.Host = viper.GetString("host")
		dbConfig.Port = viper.GetInt("port")
		dbConfig.User = viper.GetString("user")
		dbConfig.Password = viper.GetString("password")
		dbConfig.DBName = viper.GetString("dbname")
		dbConfig.SSLMode = viper.GetString("sslmode")
	}
}
