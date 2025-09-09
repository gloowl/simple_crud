package cmd

import (
	"fmt"
	"github.com/gloowl/simple_crud/src/internal/database"
	"github.com/gloowl/simple_crud/src/internal/models"
	"github.com/gloowl/simple_crud/src/internal/repository"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// var herbRepo *repository.HerbRepository

// herbCmd represents the herb command
var herbCmd = &cobra.Command{
	Use:   "herb",
	Short: "Управление травами",
	Long:  `Команды для работы с записями о лекарственных травах в базе данных.`,
	//PersistentPreRun: func(cmd *cobra.Command, args []string) {
	//	// Initialize repository
	//	db := database.GetDB()
	//  if db == nil {
	//      fmt.Println("❌ Нет подключения к БД. Проверь config/flags.")
	//      os.Exit(1)
	//  }
	//	herbRepo := repository.NewHerbRepository(db)
	//},
}

// createHerbCmd creates a new herb
var createHerbCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новую траву",
	Long:  `Создает новую запись о лекарственной траве в базе данных.`,
	Example: `  herbs-cli herb create --name "Ромашка" --latin "Matricaria chamomilla" --desc "Противовоспалительное средство"
  herbs-cli herb create --name "Белена" --latin "Hyoscyamus niger" --desc "Ядовитое растение" --poisonous`,
	RunE: createHerb,
}

// listHerbsCmd lists all herbs
var listHerbsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Показать все травы",
	Long:    `Выводит список всех лекарственных трав из базы данных.`,
	RunE:    listHerbs,
}

// getHerbCmd gets a herb by ID
var getHerbCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Получить траву по ID",
	Long:  `Выводит подробную информацию о траве с указанным ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  getHerb,
}

// updateHerbCmd updates a herb
var updateHerbCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Обновить траву",
	Long:  `Обновляет информацию о траве с указанным ID.`,
	Args:  cobra.ExactArgs(1),
	Example: `  herbs-cli herb update 1 --name "Новое название"
  herbs-cli herb update 1 --desc "Новое описание" --poisonous=false`,
	RunE: updateHerb,
}

// deleteHerbCmd deletes a herb
var deleteHerbCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Удалить траву",
	Long:  `Удаляет траву с указанным ID из базы данных.`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteHerb,
}

// searchHerbCmd searches herbs by name
var searchHerbCmd = &cobra.Command{
	Use:   "search [название]",
	Short: "Найти травы по названию",
	Long:  `Выполняет поиск трав по названию (поддерживает частичное совпадение).`,
	Args:  cobra.ExactArgs(1),
	Example: `  herbs-cli herb search "ромашка"
  herbs-cli herb search "Matricaria"`,
	RunE: searchHerbs,
}

// poisonousHerbsCmd lists poisonous herbs
var poisonousHerbsCmd = &cobra.Command{
	Use:   "poisonous",
	Short: "Показать ядовитые травы",
	Long:  `Выводит список всех ядовитых трав из базы данных.`,
	RunE:  listPoisonousHerbs,
}

func init() {
	rootCmd.AddCommand(herbCmd)

	// Add subcommands
	herbCmd.AddCommand(createHerbCmd)
	herbCmd.AddCommand(listHerbsCmd)
	herbCmd.AddCommand(getHerbCmd)
	herbCmd.AddCommand(updateHerbCmd)
	herbCmd.AddCommand(deleteHerbCmd)
	herbCmd.AddCommand(searchHerbCmd)
	herbCmd.AddCommand(poisonousHerbsCmd)

	// Flags for create command
	createHerbCmd.Flags().StringP("name", "n", "", "название травы (обязательно)")
	createHerbCmd.Flags().StringP("latin", "l", "", "латинское название")
	createHerbCmd.Flags().StringP("desc", "d", "", "описание травы")
	createHerbCmd.Flags().BoolP("poisonous", "p", false, "является ли трава ядовитой")
	createHerbCmd.Flags().StringP("image", "i", "", "путь к изображению")
	createHerbCmd.MarkFlagRequired("name")

	// Flags for update command
	updateHerbCmd.Flags().StringP("name", "n", "", "новое название травы")
	updateHerbCmd.Flags().StringP("latin", "l", "", "новое латинское название")
	updateHerbCmd.Flags().StringP("desc", "d", "", "новое описание травы")
	updateHerbCmd.Flags().BoolP("poisonous", "p", false, "является ли трава ядовитой")
	updateHerbCmd.Flags().StringP("image", "i", "", "новый путь к изображению")

	// Flags for list command
	listHerbsCmd.Flags().BoolP("table", "t", false, "вывод в табличном формате")
}

func createHerb(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	name, _ := cmd.Flags().GetString("name")
	latinName, _ := cmd.Flags().GetString("latin")
	description, _ := cmd.Flags().GetString("desc")
	isPoisonous, _ := cmd.Flags().GetBool("poisonous")
	imagePath, _ := cmd.Flags().GetString("image")

	herb := &models.Herb{
		Name:        strings.TrimSpace(name),
		LatinName:   strings.TrimSpace(latinName),
		Description: strings.TrimSpace(description),
		IsPoisonous: isPoisonous,
		ImagePath:   strings.TrimSpace(imagePath),
	}

	err := herbRepo.Create(herb)
	if err != nil {
		return fmt.Errorf("не удалось создать траву: %v", err)
	}

	fmt.Printf("✅ Трава успешно создана с ID: %d\n", herb.ID)
	fmt.Println(herb.String())
	return nil
}

func listHerbs(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	herbs, err := herbRepo.GetAll()
	if err != nil {
		return fmt.Errorf("не удалось получить список трав: %v", err)
	}

	if len(herbs) == 0 {
		fmt.Println("База данных пуста. Добавьте травы с помощью команды 'create'.")
		return nil
	}

	tableFormat, _ := cmd.Flags().GetBool("table")

	fmt.Printf("Найдено трав: %d\n\n", len(herbs))

	if tableFormat {
		// Table format
		fmt.Println(herbs[0].TableHeader())
		fmt.Println(strings.Repeat("-", 80))
		for _, herb := range herbs {
			fmt.Println(herb.TableRow())
		}
	} else {
		// Detailed format
		for i, herb := range herbs {
			if i > 0 {
				fmt.Println("\n" + strings.Repeat("-", 50))
			}
			fmt.Println(herb.String())
		}
	}

	return nil
}

func getHerb(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("неверный ID: %s", args[0])
	}

	herb, err := herbRepo.GetByID(id)
	if err != nil {
		return err
	}

	fmt.Println(herb.String())
	return nil
}

func updateHerb(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("неверный ID: %s", args[0])
	}

	// Get existing herb
	herb, err := herbRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Update fields if flags are provided
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		herb.Name = strings.TrimSpace(name)
	}
	if cmd.Flags().Changed("latin") {
		latin, _ := cmd.Flags().GetString("latin")
		herb.LatinName = strings.TrimSpace(latin)
	}
	if cmd.Flags().Changed("desc") {
		desc, _ := cmd.Flags().GetString("desc")
		herb.Description = strings.TrimSpace(desc)
	}
	if cmd.Flags().Changed("poisonous") {
		poisonous, _ := cmd.Flags().GetBool("poisonous")
		herb.IsPoisonous = poisonous
	}
	if cmd.Flags().Changed("image") {
		image, _ := cmd.Flags().GetString("image")
		herb.ImagePath = strings.TrimSpace(image)
	}

	err = herbRepo.Update(herb)
	if err != nil {
		return fmt.Errorf("не удалось обновить траву: %v", err)
	}

	fmt.Printf("✅ Трава с ID %d успешно обновлена\n", herb.ID)
	fmt.Println(herb.String())
	return nil
}

func deleteHerb(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("неверный ID: %s", args[0])
	}

	// Show herb before deletion
	herb, err := herbRepo.GetByID(id)
	if err != nil {
		return err
	}

	fmt.Println("Удаляется следующая трава:")
	fmt.Println(herb.String())
	fmt.Print("\nВы уверены? (y/N): ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "y" && confirmation != "Y" {
		fmt.Println("Удаление отменено.")
		return nil
	}

	err = herbRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("не удалось удалить траву: %v", err)
	}

	fmt.Printf("✅ Трава с ID %d успешно удалена\n", id)
	return nil
}

func searchHerbs(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	searchTerm := args[0]
	herbs, err := herbRepo.Search(searchTerm)
	if err != nil {
		return fmt.Errorf("ошибка поиска: %v", err)
	}

	if len(herbs) == 0 {
		fmt.Printf("Травы с названием '%s' не найдены.\n", searchTerm)
		return nil
	}

	fmt.Printf("Найдено трав по запросу '%s': %d\n\n", searchTerm, len(herbs))

	for i, herb := range herbs {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("-", 50))
		}
		fmt.Println(herb.String())
	}

	return nil
}

func listPoisonousHerbs(cmd *cobra.Command, args []string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("❌ нет соединения с БД (database.GetDB() == nil)")
	}
	herbRepo := repository.NewHerbRepository(db)

	herbs, err := herbRepo.GetPoisonous()
	if err != nil {
		return fmt.Errorf("не удалось получить список ядовитых трав: %v", err)
	}

	if len(herbs) == 0 {
		fmt.Println("В базе данных нет записей о ядовитых травах.")
		return nil
	}

	fmt.Printf("⚠️  Найдено ядовитых трав: %d\n\n", len(herbs))

	for i, herb := range herbs {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("-", 50))
		}
		fmt.Println(herb.String())
	}

	return nil
}
