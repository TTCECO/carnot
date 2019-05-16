// Copyright 2019 The TTC Authors
// This file is part of the TTC library.
//
// The TTC library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The TTC library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the TTC library. If not, see <http://www.gnu.org/licenses/>.

package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

const (
	orderID = "orderID"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	//r.HandleFunc("/%s/tcchan", namesHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/tcchan/order/{%s}", orderID), resolveOrderHandler(cdc, cliCtx, storeName)).Methods("GET")
	//r.HandleFunc(fmt.Sprintf("/%s/tcchan", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	//r.HandleFunc(fmt.Sprintf("/%s/tcchan/{%s}", storeName, restName), resolveNameHandler(cdc, cliCtx, storeName)).Methods("GET")
	//r.HandleFunc(fmt.Sprintf("/%s/tcchan/{%s}/whois", storeName, restName), whoIsHandler(cdc, cliCtx, storeName)).Methods("GET")
}

func resolveOrderHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[orderID]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/order/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
