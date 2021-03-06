package backup

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/Percona-Lab/percona-xtradb-cluster-operator/pkg/apis/pxc/v1alpha1"
)

// Job returns the backup job
func Job(cr *api.PerconaXtraDBBackup) batchv1.Job {
	pvc := corev1.Volume{
		Name: "backup",
	}
	pvsource := corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: "backup-volume",
	}
	pvc.PersistentVolumeClaim = &pvsource

	return batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "xtrabackup-job",
			Namespace: cr.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "xtrabackup",
							Image:   "perconalab/backupjob-openshift",
							Command: []string{"bash", "/usr/bin/backup.sh"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "backup",
									MountPath: "/backup",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "NODE_NAME",
									Value: "cluster1-pxc-nodes",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						pvc,
					},
				},
			},
			BackoffLimit: func(i int32) *int32 { return &i }(4),
		},
	}
}
