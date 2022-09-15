// Code generated by smithy-go-codegen DO NOT EDIT.

package ecrpublic

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Informs Amazon ECR that the image layer upload has completed for a specified
// public registry, repository name, and upload ID. You can optionally provide a
// sha256 digest of the image layer for data validation purposes. When an image is
// pushed, the CompleteLayerUpload API is called once per each new image layer to
// verify that the upload has completed. This operation is used by the Amazon ECR
// proxy and is not generally used by customers for pulling and pushing images. In
// most cases, you should use the docker CLI to pull, tag, and push images.
func (c *Client) CompleteLayerUpload(ctx context.Context, params *CompleteLayerUploadInput, optFns ...func(*Options)) (*CompleteLayerUploadOutput, error) {
	if params == nil {
		params = &CompleteLayerUploadInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CompleteLayerUpload", params, optFns, c.addOperationCompleteLayerUploadMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CompleteLayerUploadOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CompleteLayerUploadInput struct {

	// The sha256 digest of the image layer.
	//
	// This member is required.
	LayerDigests []string

	// The name of the repository in a public registry to associate with the image
	// layer.
	//
	// This member is required.
	RepositoryName *string

	// The upload ID from a previous InitiateLayerUpload operation to associate with
	// the image layer.
	//
	// This member is required.
	UploadId *string

	// The AWS account ID associated with the registry to which to upload layers. If
	// you do not specify a registry, the default public registry is assumed.
	RegistryId *string

	noSmithyDocumentSerde
}

type CompleteLayerUploadOutput struct {

	// The sha256 digest of the image layer.
	LayerDigest *string

	// The public registry ID associated with the request.
	RegistryId *string

	// The repository name associated with the request.
	RepositoryName *string

	// The upload ID associated with the layer.
	UploadId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCompleteLayerUploadMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpCompleteLayerUpload{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpCompleteLayerUpload{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpCompleteLayerUploadValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCompleteLayerUpload(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opCompleteLayerUpload(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ecr-public",
		OperationName: "CompleteLayerUpload",
	}
}
