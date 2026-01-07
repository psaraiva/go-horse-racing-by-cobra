package main

import (
	"os"

	"github.com/psaraiva/go-horse-racing-by-cobra/internal"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lab-go-horse-racing-by-cobra",
	Short: "The GO race of the horses",
	Long: `Horses run free on the hills...
Horses race against other horses under human supervision...`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("horse-label", internal.HorseLabelDefault, "Horse label, only one char")
	rootCmd.Flags().Int("horses-quantity", internal.HorseQuantityDefault, "Horses quantity, 2 to 99")
	rootCmd.Flags().Int("score-target", internal.ScoreTargetDefault, "Score target, 15 to 100")
	rootCmd.Flags().String("game-timeout", internal.GameTimeoutDefault, "Game timeout, 10s to 90s")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		horseLabel, _ := rootCmd.Flags().GetString("horse-label")
		horsesQuantity, _ := rootCmd.Flags().GetInt("horses-quantity")
		scoreTarget, _ := rootCmd.Flags().GetInt("score-target")
		gameTimeout, _ := rootCmd.Flags().GetString("game-timeout")

		internal.Run(internal.Input{
			HorseLabel:     horseLabel,
			HorsesQuantity: horsesQuantity,
			ScoreTarget:    scoreTarget,
			GameTimeout:    gameTimeout,
		})
	}
}
