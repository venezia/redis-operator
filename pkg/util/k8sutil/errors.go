package k8sutil

import (
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

func IsResourceAlreadyExistsError(err error) bool {
	return k8serrors.IsAlreadyExists(err)
}

func IsResourceNotFoundError(err error) bool {
	return k8serrors.IsNotFound(err)
}
