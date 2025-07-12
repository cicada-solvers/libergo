package titler

import (
	"config"
	"fmt"
)

// PrintTitle displays a decorative title banner if the configuration does not hide the title.
func PrintTitle(title string) {
	configuration, _ := config.LoadConfig()
	if configuration.HideTitle {
		return
	}

	fmt.Println("                                                                                                              ")
	fmt.Println("                                                                                                              ")
	fmt.Println("                                                                                                              ")
	fmt.Println("   ########        ###        ########            #   %%                                                      ")
	fmt.Println(" %##  # ####%%%%# %###  #%##% ##%  ##%#  %%       #   %%                                                      ")
	fmt.Println("   ##############%## #####%%%##%##%%#    %%       %#  %%#%%%#     #%%%#   %# %%#  #%%%##%    #%%%#            ")
	fmt.Println("      ## ###     #% #       %#####       %%       @#  %@    %@  #@#   %%  %@#    %%    @%  #@#   %@           ")
	fmt.Println("        %###      ####      ##%%         %%       @#  %%    #@# %@@@@@@@  %%     @#    %%  @%     @%          ")
	fmt.Println("          #####   %## ##  #%%#           %%       @#  %%    #@  %%        %%     @#    %%  %%     @#          ")
	fmt.Println("                  ##                     %@#####  @#  %%@##%@#   %@%##%#  %%     #@%##%%%   %@##%@#           ")
	fmt.Println("                   ##                                                                  %%                     ")
	fmt.Println("                                                                                 %####%@#                     ")
	fmt.Println("                                                                                                              ")
	fmt.Println("==============================================================================================================")
	fmt.Println(title)
	fmt.Println("==============================================================================================================")
}
