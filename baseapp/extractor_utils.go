package baseapp

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"os"
)

//Denominations represents the list of all denominations at a block height
type Denominations = []string

var genesisDenominations = map[string]Denominations{
	//TODO : Parse genesis denominations from genesis.json 
	// https://github.com/ChorusOne/hipparchus/issues/78

	"columbus-1": []string{"ukrw", "uluna", "umnt", "usdr", "uusd"},
	"columbus-2": []string{"ukrw", "uluna", "umnt", "usdr", "uusd"},
	"columbus-3": []string{"ukrw", "uluna", "umnt", "usdr", "uusd"},

	"cosmoshub-1": []string{"uatom"},
	"cosmoshub-2": []string{"uatom"},
	"cosmoshub-3": []string{"uatom"},
}

func copyFile(destination string, source string) error {
	sourceFile, err := os.OpenFile(source, os.O_RDONLY, 0644)
	if err != nil {
		// It is not an error if source file does not exists
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error: (%v) while trying to open source file while copying", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error: (%v) while trying to open destination file while copying", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("error: (%v) while trying to copy source file into destination", err)
	}

	return nil
}

func commitUncheckedFiles(ctx sdk.Context) {
	for _, key := range []string{"delegations", "unbond", "balance", "rewards", "commission"} {
		err := copyFile(fmt.Sprintf("./extract/progress/%s.%d.%s", key, ctx.BlockHeight(), ctx.ChainID()), fmt.Sprintf("./extract/unchecked/%s.%d.%s", key, ctx.BlockHeight(), ctx.ChainID()))
		if err != nil {
			panic(fmt.Sprintf("error: (%v) while commiting unchecked file\n", err))
		}
		// No need for the file now
		if err := os.Remove(fmt.Sprintf("./extract/unchecked/%s.%d.%s", key, ctx.BlockHeight(), ctx.ChainID())); err != nil && !os.IsNotExist(err) {
			panic(fmt.Sprintf("error: (%v) while removing unchecked file after commiting\n", err))
		}
	}
}

func deleteUncheckedFiles(ctx sdk.Context) {
	for _, key := range []string{"delegations", "unbond", "balance", "rewards", "commission"} {
		if err := os.Remove(fmt.Sprintf("./extract/unchecked/%s.%d.%s", key, ctx.BlockHeight(), ctx.ChainID())); err != nil && !os.IsNotExist(err) {
			panic(fmt.Sprintf("error: (%v) while removing unchecked file\n", err))
		}
	}
}

func (app *BaseApp) SetExtractDataMode() {
	app.extractData = true
}

func (app *BaseApp) GetExtractDataMode() bool {
	return app.extractData
}
