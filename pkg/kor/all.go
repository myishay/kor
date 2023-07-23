package kor

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"os"
)

type ResourceDiff struct {
	resourceType string
	diff         []string
}

func getUnusedCMs(kubeClient *kubernetes.Clientset, namespace string) ResourceDiff {
	cmDiff, err := processNamespaceCM(kubeClient, namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get %s namespace %s: %v\n", "configmaps", namespace, err)
	}
	namespaceCMDiff := ResourceDiff{"ConfigMap", cmDiff}
	return namespaceCMDiff
}

func getUnusedSVCs(kubeClient *kubernetes.Clientset, namespace string) ResourceDiff {
	svcDiff, err := ProcessNamespaceServices(kubeClient, namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get %s namespace %s: %v\n", "services", namespace, err)
	}
	namespaceSVCDiff := ResourceDiff{"Service", svcDiff}
	return namespaceSVCDiff
}

func getUnusedSecrets(kubeClient *kubernetes.Clientset, namespace string) ResourceDiff {
	secretDiff, err := processNamespaceSecret(kubeClient, namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get %s namespace %s: %v\n", "secrets", namespace, err)
	}
	namespaceSecretDiff := ResourceDiff{"Secret", secretDiff}
	return namespaceSecretDiff
}

func getUnusedServiceAccounts(kubeClient *kubernetes.Clientset, namespace string) ResourceDiff {
	saDiff, err := processNamespaceSA(kubeClient, namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get %s namespace %s: %v\n", "serviceaccounts", namespace, err)
	}
	namespaceSADiff := ResourceDiff{"ServiceAccount", saDiff}
	return namespaceSADiff
}

func GetUnusedAll(namespace string) {
	var kubeClient *kubernetes.Clientset
	var namespaces []string
	var allDiffs []ResourceDiff

	kubeClient = GetKubeClient()

	namespaces = SetNamespaceList(namespace, kubeClient)
	for _, namespace := range namespaces {
		namespaceCMDiff := getUnusedCMs(kubeClient, namespace)
		allDiffs = append(allDiffs, namespaceCMDiff)
		namespaceSVCDiff := getUnusedSVCs(kubeClient, namespace)
		allDiffs = append(allDiffs, namespaceSVCDiff)
		namespaceSecretDiff := getUnusedSecrets(kubeClient, namespace)
		allDiffs = append(allDiffs, namespaceSecretDiff)
		namespaceSADiff := getUnusedServiceAccounts(kubeClient, namespace)
		allDiffs = append(allDiffs, namespaceSADiff)
		output := FormatOutputAll(namespace, allDiffs)
		fmt.Println(output)
		fmt.Println()
	}
}