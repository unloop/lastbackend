//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package provider

import (
	"github.com/lastbackend/lastbackend/internal/master/state/job/provider"
	"github.com/lastbackend/lastbackend/internal/master/state/job/provider/http"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

func New(specProvider models.JobSpecProvider) (provider.JobProvider, error) {

	if specProvider.Http != nil {
		if specProvider.Http.Endpoint != models.EmptyString {
			return http.New(specProvider.Http)
		}
	}

	return nil, nil
}
