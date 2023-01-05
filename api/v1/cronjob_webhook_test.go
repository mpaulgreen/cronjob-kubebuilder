/** ginkgo best practices are broken (tossed out of window). DO NOT COPY**/

package v1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	testNameSpace         = "default"
	testCorrectName       = "cronjob-sample"
	testCorrectSchedule   = "*/1 * * * *"
	tesIncorrectName      = "This is a very long string greater than fifty two charters repeat This is a very long string greater than fifty two charters repeat This is a very long string greater than fifty two charters "
	testIncorrectSchedule = "*/"
)

var _ = Describe("For Validating webhook of cronjob", func() {
	Context("Creation of a cron job", func() {

		cronjobLookupKey := types.NamespacedName{Name: testCorrectName, Namespace: testNameSpace}
		createdCronjob := &CronJob{}
		correctCronJob := &CronJob{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "batch.tutorial.kubebuilder.io/v1",
				Kind:       "CronJob",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      testCorrectName,
				Namespace: testNameSpace,
			},
			Spec: CronJobSpec{
				Schedule: testCorrectSchedule,
				JobTemplate: batchv1.JobTemplateSpec{
					Spec: batchv1.JobSpec{
						// For simplicity, we only fill out the required fields.
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								// For simplicity, we only fill out the required fields.
								Containers: []v1.Container{
									{
										Name:  "test-container",
										Image: "test-image",
									},
								},
								RestartPolicy: v1.RestartPolicyOnFailure,
							},
						},
					},
				},
			},
		}

		incorrectCronJob := &CronJob{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "batch.tutorial.kubebuilder.io/v1",
				Kind:       "CronJob",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      tesIncorrectName,
				Namespace: testNameSpace,
			},
			Spec: CronJobSpec{
				Schedule: testCorrectSchedule,
				JobTemplate: batchv1.JobTemplateSpec{
					Spec: batchv1.JobSpec{
						// For simplicity, we only fill out the required fields.
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								// For simplicity, we only fill out the required fields.
								Containers: []v1.Container{
									{
										Name:  "test-container",
										Image: "test-image",
									},
								},
								RestartPolicy: v1.RestartPolicyOnFailure,
							},
						},
					},
				},
			},
		}

		It("With correct name and schedule format ", func() {
			Expect(k8sClient.Create(ctx, correctCronJob)).Should(Succeed())
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("With incorrect name and correct schedule format ", func() {
			err := k8sClient.Create(ctx, incorrectCronJob)
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError("CronJob.batch.tutorial.kubebuilder.io \"This is a very long string greater than fifty two charters repeat This is a very long string greater than fifty two charters repeat This is a very long string greater than fifty two charters \" is invalid: metadata.name: Invalid value: \"This is a very long string greater than fifty two charters repeat This is a very long string greater than fifty two charters repeat This is a very long string greater than fifty two charters \": a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character (e.g. 'example.com', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*')"))
		})

		It("With correct name and incorrect schedule format ", func() {

			correctCronJob.Spec.Schedule = testIncorrectSchedule
			err := k8sClient.Create(ctx, correctCronJob)
			Expect(err).ShouldNot(BeNil())
			println("Mriganka>>>>>>>>>>>>>>>", err)
			Expect(err).Should(MatchError("admission webhook \"vcronjob.kb.io\" denied the request: CronJob.batch.tutorial.kubebuilder.io \"cronjob-sample\" is invalid: spec.schedule: Invalid value: \"*/\": Expected exactly 5 fields, found 1: */"))
		})
	})

	Context("Upation of cron job", func() {
		var cronjobLookupKey types.NamespacedName
		var createdCronjob, correctCronJob *CronJob
		BeforeEach(func() {
			cronjobLookupKey = types.NamespacedName{Name: testCorrectName, Namespace: testNameSpace}
			createdCronjob = &CronJob{}
			correctCronJob = &CronJob{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch.tutorial.kubebuilder.io/v1",
					Kind:       "CronJob",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample",
					Namespace: testNameSpace,
				},
				Spec: CronJobSpec{
					Schedule: "*/1 1 * * *",
					JobTemplate: batchv1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							// For simplicity, we only fill out the required fields.
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									// For simplicity, we only fill out the required fields.
									Containers: []v1.Container{
										{
											Name:  "test-container",
											Image: "test-image",
										},
									},
									RestartPolicy: v1.RestartPolicyOnFailure,
								},
							},
						},
					},
				},
			}
		})

		It("With correct name and schedule format ", func() {
			Expect(k8sClient.Create(ctx, correctCronJob)).Should(Succeed())
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			print(createdCronjob)
			createdCronjob.Name = testCorrectName
			createdCronjob.Spec.Schedule = testCorrectSchedule
			err := k8sClient.Update(ctx, createdCronjob)
			Expect(err).Should(BeNil())
			err = k8sClient.Delete(ctx, correctCronJob)
			Expect(err).Should(BeNil())

		})

		It("With correct name and incorrect schedule format ", func() {
			Expect(k8sClient.Create(ctx, correctCronJob)).Should(Succeed())
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			createdCronjob.Spec.Schedule = testIncorrectSchedule
			err := k8sClient.Update(ctx, createdCronjob)
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError("admission webhook \"vcronjob.kb.io\" denied the request: CronJob.batch.tutorial.kubebuilder.io \"cronjob-sample\" is invalid: spec.schedule: Invalid value: \"*/\": Expected exactly 5 fields, found 1: */"))
		})

	})
})
