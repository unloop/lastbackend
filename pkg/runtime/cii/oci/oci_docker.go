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

package oci

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	stg "github.com/containers/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/pkg/errors"
)

const (
	DefaultLocalRegistry     = "localhost"
	DefaultLRootPath         = "/var/lib/lastbackend/"
	DefaultTransport         = "docker://"
	minimumTruncatedIDLength = 3
)

var localRegistryPrefix = DefaultLocalRegistry + "/"

type Docker struct {
	unqualifiedSearchRegistries []string
	insecureRegistryCIDRs       []*net.IPNet
	indexConfigs                map[string]*indexInfo
	defaultTransport            string
	config                      ConfigOci
	store                       stg.Store
}

type indexInfo struct {
	name   string
	secure bool
}

func NewOci(config ConfigOci) (*Docker, error) {

	options, err := stg.DefaultStoreOptions(false, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create default image store options")
	}

	options.RunRoot = DefaultLRootPath
	options.GraphRoot = DefaultLRootPath

	if config.RunRoot != "" {
		options.RunRoot = config.RunRoot
		options.GraphRoot = config.Root
	}

	if config.StorageDriver != "" {
		options.GraphDriverName = config.StorageDriver
	}

	store, err := stg.GetStore(options)
	if err != nil {
		return nil, err
	}

	d := new(Docker)
	d.config = config
	d.store = store
	d.defaultTransport = DefaultTransport

	return d, nil
}

func (d *Docker) Auth(ctx context.Context, secret *models.SecretAuthData) (string, error) {
	config := models.AuthConfig{
		Username: secret.Username,
		Password: secret.Password,
	}
	js, err := json.Marshal(config)
	if err != nil {
		return models.EmptyString, err
	}
	return base64.URLEncoding.EncodeToString(js), nil
}

func (d *Docker) Pull(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error) {

	imageName := spec.Name

	systemCtx := &types.SystemContext{}
	options := &copy.Options{
		ReportWriter: out,
	}

	srcSystemContext, srcRef, err := d.prepareReference(systemCtx, imageName)
	if err != nil {
		return nil, err
	}
	options.SourceCtx = srcSystemContext

	policy, err := signature.DefaultPolicy(systemCtx)
	if err != nil {
		return nil, err
	}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return nil, err
	}

	dest := imageName
	if srcRef.DockerReference() != nil {
		dest = srcRef.DockerReference().String()
	}
	destRef, err := storage.Transport.ParseStoreReference(d.store, dest)
	if err != nil {
		return nil, err
	}
	_, err = copy.Image(ctx, policyContext, destRef, srcRef, options)
	if err != nil {
		return nil, err
	}

	image, err := storage.Transport.GetStoreImage(d.store, destRef)
	if err != nil {
		return nil, err
	}

	img := new(models.Image)
	img.Meta.ID = image.ID
	img.Meta.Name = destRef.DockerReference().Name()
	img.Meta.Digest = image.Digest.String()

	return img, nil
}

func (d *Docker) Remove(ctx context.Context, image string) error {
	return nil
}

func (d *Docker) Push(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error) {
	return nil, nil
}

func (d *Docker) Build(ctx context.Context, stream io.Reader, spec *models.SpecBuildImage, out io.Writer) (*models.Image, error) {
	return nil, nil
}

func (d *Docker) List(ctx context.Context, filters ...string) ([]*models.Image, error) {
	return nil, nil
}

func (d *Docker) Inspect(ctx context.Context, nameOrID string) (*models.Image, error) {
	ref, err := d.getRef(nameOrID)
	if err != nil {
		return nil, err
	}
	image, err := storage.Transport.GetStoreImage(d.store, ref)
	if err != nil {
		return nil, err
	}

	img := new(models.Image)
	img.Meta.ID = image.ID
	img.Meta.Name = ref.DockerReference().Name()
	img.Meta.Digest = image.Digest.String()

	return img, nil
}

func (d *Docker) Subscribe(ctx context.Context) (chan *models.Image, error) {
	return nil, nil
}

func (d *Docker) Close() error {
	return nil
}

func (d *Docker) prepareReference(inputSystemContext *types.SystemContext, imageName string) (*types.SystemContext, types.ImageReference, error) {
	srcRef, err := d.remoteImageReference(imageName)
	if err != nil {
		return nil, nil, err
	}

	sc := types.SystemContext{}
	if inputSystemContext != nil {
		sc = *inputSystemContext // A shallow copy
	}

	if srcRef.DockerReference() != nil {
		hostname := reference.Domain(srcRef.DockerReference())
		if secure := d.isSecureIndex(hostname); !secure {
			sc.DockerInsecureSkipTLSVerify = types.OptionalBoolTrue
		}
	}

	return &sc, srcRef, nil
}

func (d *Docker) remoteImageReference(imageName string) (types.ImageReference, error) {
	if imageName == "" {
		return nil, stg.ErrNotAnImage
	}

	srcRef, err := alltransports.ParseImageName(imageName)
	if err != nil {
		if d.defaultTransport == "" {
			return nil, err
		}
		srcRef2, err2 := alltransports.ParseImageName(d.defaultTransport + imageName)
		if err2 != nil {
			return nil, err
		}
		srcRef = srcRef2
	}

	return srcRef, nil
}

func (d *Docker) isSecureIndex(indexName string) bool {
	if index, ok := d.indexConfigs[indexName]; ok {
		return index.secure
	}

	host, _, err := net.SplitHostPort(indexName)
	if err != nil {
		// assume indexName is of the form `host` without the port and go on.
		host = indexName
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		ip := net.ParseIP(host)
		if ip != nil {
			addrs = []net.IP{ip}
		}

		// if ip == nil, then `host` is neither an IP nor it could be looked up,
		// either because the index is unreachable, or because the index is behind an HTTP proxy.
		// So, len(addrs) == 0 and we're not aborting.
	}

	// Try CIDR notation only if addrs has any elements, i.e. if `host`'s IP could be determined.
	for _, addr := range addrs {
		for _, ipnet := range d.insecureRegistryCIDRs {
			// check if the addr falls in the subnet
			if ipnet.Contains(addr) {
				return false
			}
		}
	}

	return true
}

func (d *Docker) getRef(name string) (types.ImageReference, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		ref2, err2 := storage.Transport.ParseStoreReference(d.store, "@"+name)
		if err2 != nil {
			ref3, err3 := storage.Transport.ParseStoreReference(d.store, name)
			if err3 != nil {
				return nil, err
			}
			ref2 = ref3
		}
		ref = ref2
	}
	return ref, nil
}
