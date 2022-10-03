package cmd

import (
	"log"
	"os"
	"time"

	"github.com/siddmoitra/btech-minicash-2-poc/backend/handler"
	"github.com/siddmoitra/btech-minicash-2-poc/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "btech-minicash-2-poc",
	Short: "A MiniCash 2.0 futuristic PoC",
	Long: `
	MiniCash 2.0 is PoC to verify the hypothesis of segregating products from the core LOS/LMS system.

	This also includes the high-level structure of how
	a MiniCash 2.0 - a Self Contained System interacts with
	various other systems within the retail client.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var shopCmd = &cobra.Command{
	Use:   "shop",
	Short: "Start a shopping session on retail e-Commerce store",
	Long:  "This demo showcases the end-consumer journey when using digital shop and chooses to buy the product(s) using the pay by MiniCash 2.0",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Start end-consumer journey on retailer's digital e-commerce website or through a retailer's in-shop system assisted by sales agent ...")

		handler.ECommerceLoadData()
		handler.MiniCashLoadData()

		salesOrder, err := ECommerceStartShopping()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		utils.LineBreak(2)
		utils.SpinnerMessage("  You are now redirected to MiniCash 2.0 in 5 seconds", time.Second*5)
		utils.LineBreak(2)

		paymentResponse, err := MiniCashInitiatePayment(salesOrder)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		utils.LineBreak(2)
		utils.SpinnerMessage("  You are now back in E-Commerce shop with MiniCash 2.0 payment response", time.Second*2)
		utils.LineBreak(2)

		ECommerceHandlePaymentStatus(salesOrder, paymentResponse)

		handler.MiniCashWriteData()
	},
}

var walletStatusCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Show details of e-Wallet of a customer in MiniCash 2.0",
	Long:  "This demo showcases the e-Wallet feature of the MiniCash 2.0 from the end-user point of view",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Start end-consumer journey on MiniCash 2.0 using omnichannel frontend ...")

		handler.ECommerceLoadData()
		handler.MiniCashLoadData()
		MiniCashCustomerWallet()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.btech-minicash-2-poc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(shopCmd)
	rootCmd.AddCommand(walletStatusCmd)
}
