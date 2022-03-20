package broadcasters

import "testing"

func TestSlack(t *testing.T) {
	msg := `[CRITICAL]: Resources->
	check this out
	namespace: kube-system, kind: DaemonSet, name: kube-proxy
	status:failed()
	Control: C-0014 Access Kubernetes dashboard
	Control: C-0057 Privileged container
	Fix: :
	Remove: spec.template.spec.containers[0].securityContext.privileged
	Control: C-0045 Writable hostPath mount
	Control: C-0020 Mount service principal
	Fix: :
	Remove: spec.template.spec.volumes[0].hostPath.path
	Control: C-0012 Applications credentials in configuration files
	Control: C-0048 HostPath mount
	Fix: :
	Remove: spec.template.spec.volumes[0].hostPath.path
	original yaml:
	apiVersion: apps/v1
	kind: DaemonSet
	metadata:
	 annotations:
	   deprecated.daemonset.template.generation: "8"
	 creationTimestamp: "2021-07-07T05:45:02Z"
	 generation: 8
	 labels:
	   addonmanager.kubernetes.io/mode: Reconcile
	   k8s-app: kube-proxy
	 name: kube-proxy
	 namespace: kube-system
	 resourceVersion: "94208546"
	 uid: 8f2390ae-0019-4e6a-96cf-0bdfe15b4be0
	spec:
	 revisionHistoryLimit: 10
	 selector:
	   matchLabels:
		 k8s-app: kube-proxy
	 template:
	   metadata:
		 creationTimestamp: null
		 labels:
		   k8s-app: kube-proxy
	   spec:
		 containers:
		 - command:
		   - /bin/sh
		   - -c
		   - kube-proxy --cluster-cidr=10.36.0.0/14 --oom-score-adj=-998 --v=2 --feature-gates=DynamicKubeletConfig=false,InTreePluginAWSUnregister=true,InTreePluginAzureDiskUnregister=true,InTreePluginOpenStackUnregister=true,InTreePluginvSphereUnregister=true,PodDeletionCost=true
			 --iptables-sync-period=1m --iptables-min-sync-period=10s --ipvs-sync-period=1m
			 --ipvs-min-sync-period=10s --detect-local-mode=NodeCIDR 1>>/var/log/kube-proxy.log
			 2>&1
		   env:
		   - name: KUBERNETES_SERVICE_HOST
			 value: XXXXXX
		   image: gke.gcr.io/kube-proxy-amd64:v1.21.6-gke.1500
		   imagePullPolicy: IfNotPresent
		   name: kube-proxy
		   resources:
			 requests:
			   cpu: 100m
		   securityContext:
			 privileged: true
		   terminationMessagePath: /dev/termination-log
		   terminationMessagePolicy: File
		   volumeMounts:
		   - mountPath: /var/log
			 name: varlog
		   - mountPath: /run/xtables.lock
			 name: xtables-lock
		   - mountPath: /lib/modules
			 name: lib-modules
			 readOnly: true
		 dnsPolicy: ClusterFirst
		 hostNetwork: true
		 nodeSelector:
		   kubernetes.io/os: linux
		   node.kubernetes.io/kube-proxy-ds-ready: "true"
		 priorityClassName: system-node-critical
		 restartPolicy: Always
		 schedulerName: default-scheduler
		 securityContext: {}
		 serviceAccount: kube-proxy
		 serviceAccountName: kube-proxy
		 terminationGracePeriodSeconds: 30
		 tolerations:
		 - effect: NoExecute
		   operator: Exists
		 - effect: NoSchedule
		   operator: Exists
		 volumes:
		 - hostPath:
			 path: /var/log
			 type: ""
		   name: varlog
		 - hostPath:
			 path: /run/xtables.lock
			 type: FileOrCreate
		   name: xtables-lock
		 - hostPath:
			 path: /lib/modules
			 type: ""
		   name: lib-modules
	 updateStrategy:
	   rollingUpdate:
		 maxSurge: 0
		 maxUnavailable: 10%
	   type: RollingUpdate
	==============
	12:47
	[CRITICAL]: Resources->
	check this out
	namespace: , kind: User, name: system:konnectivity-server
	status:passed()
	Control: C-0014 Access Kubernetes dashboard
	Control: C-0002 Exec into container
	Control: C-0015 List Kubernetes secrets
	Control: C-0031 Delete Kubernetes events
	Control: C-0007 Data Destruction
	Control: C-0035 Cluster-admin binding
	Control: C-0037 CoreDNS poisoning
	original yaml:
	apiGroup: rbac.authorization.k8s.io
	kind: User
	name: system:konnectivity-server
	relatedObjects:
	- apiVersion: rbac.authorization.k8s.io/v1
	 kind: ClusterRoleBinding
	 metadata:
	   annotations:
		 components.gke.io/component-name: konnectivitynetworkproxy-combined
		 components.gke.io/component-version: 1.3.2
		 components.gke.io/layer: addon
	   creationTimestamp: "2021-11-02T02:33:45Z"
	   labels:
		 addonmanager.kubernetes.io/mode: Reconcile
	   name: system:konnectivity-server
	   resourceVersion: "55353540"
	   uid: d017a8c1-9fbd-4327-930c-ab699418cb30
	 roleRef:
	   apiGroup: rbac.authorization.k8s.io
	   kind: ClusterRole
	   name: system:auth-delegator
	 subjects:
	 - apiGroup: rbac.authorization.k8s.io
	   kind: User
	   name: system:konnectivity-server
	- apiVersion: rbac.authorization.k8s.io/v1
	 kind: ClusterRole
	 metadata:
	   annotations:
		 rbac.authorization.kubernetes.io/autoupdate: "true"
	   creationTimestamp: "2021-07-07T05:44:34Z"
	   labels:
		 kubernetes.io/bootstrapping: rbac-defaults
	   name: system:auth-delegator
	   resourceVersion: "57"
	   uid: efcb51cd-c62c-44ed-8575-35c241b7595a
	 rules:
	 - apiGroups:
	   - authentication.k8s.io
	   resources:
	   - tokenreviews
	   verbs:
	   - create
	 - apiGroups:
	   - authorization.k8s.io
	   resources:
	   - subjectaccessreviews
	   verbs:
	   - create
	==============
	12:48
	[CRITICAL]: Resources->
	check this out
	namespace: armo-system, kind: Deployment, name: armo-notification-service
	status:passed()
	Control: C-0014 Access Kubernetes dashboard
	Control: C-0057 Privileged container
	Control: C-0045 Writable hostPath mount
	Control: C-0020 Mount service principal
	Control: C-0012 Applications credentials in configuration files
	Control: C-0048 HostPath mount
	original yaml:
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	 annotations:
	   deployment.kubernetes.io/revision: "8"
	   meta.helm.sh/release-name: armo
	   meta.helm.sh/release-namespace: armo-system
	 creationTimestamp: "2022-01-17T08:42:45Z"
	 generation: 8
	 labels:
	   app: armo-notification-service
	   app.kubernetes.io/managed-by: Helm
	   helm.sh/chart: armo-cluster-components-1.6.7
	   tier: armo-system-control-plane
	 name: armo-notification-service
	 namespace: armo-system
	 resourceVersion: "112927117"
	 uid: 042acd08-1bff-40d8-a6d9-50e388558bcb
	spec:
	 progressDeadlineSeconds: 600
	 replicas: 1
	 revisionHistoryLimit: 10
	 selector:
	   matchLabels:
		 app.kubernetes.io/instance: armo
		 app.kubernetes.io/name: armo-notification-service
		 tier: armo-system-control-plane
	 strategy:
	   rollingUpdate:
		 maxSurge: 25%
		 maxUnavailable: 25%
	   type: RollingUpdate
	 template:
	   metadata:
		 creationTimestamp: null
		 labels:
		   app: armo-notification-service
		   app.kubernetes.io/instance: armo
		   app.kubernetes.io/name: armo-notification-service
		   helm.sh/chart: armo-cluster-components-1.6.7
		   helm.sh/revision: "48"
		   tier: armo-system-control-plane
	   spec:
		 automountServiceAccountToken: false
		 containers:
		 - args:
		   - -alsologtostderr
		   - -v=4
		   - 2>&1
		   env:
		   - name: CA_CUSTOMER_GUID
			 value: XXXXXX
			 valueFrom:
			   configMapKeyRef:
				 key: accountGuid
				 name: armo-be-config
		   - name: CA_CLUSTER_NAME
			 value: XXXXXX
			 valueFrom:
			   configMapKeyRef:
				 key: clusterName
				 name: armo-be-config
		   - name: MASTER_NOTIFICATION_SERVER_HOST
			 value: XXXXXX
			 valueFrom:
			   configMapKeyRef:
				 key: masterNotificationServer
				 name: armo-be-config
		   - name: MASTER_NOTIFICATION_SERVER_ATTRIBUTES
			 value: XXXXXX
		   - name: CA_NOTIFICATION_SERVER_WS_PORT
			 value: XXXXXX
		   - name: CA_NOTIFICATION_SERVER_PORT
			 value: XXXXXX
		   image: quay.io/armosec/notification-server-ubi:89
		   imagePullPolicy: Always
		   name: armo-notification-service
		   ports:
		   - containerPort: 8001
			 name: websocket
			 protocol: TCP
		   - containerPort: 8002
			 name: rest-api
			 protocol: TCP
		   resources:
			 limits:
			   cpu: 100m
			   memory: 50Mi
			 requests:
			   cpu: 10m
			   memory: 10Mi
		   terminationMessagePath: /dev/termination-log
		   terminationMessagePolicy: File
		 dnsPolicy: ClusterFirst
		 restartPolicy: Always
		 schedulerName: default-scheduler
		 securityContext: {}
		 terminationGracePeriodSeconds: 30
		 volumes:
		 - configMap:
			 defaultMode: 420
			 items:
			 - key: clusterData
			   path: clusterData.json
			 name: armo-be-config
		   name: armo-be-config
	==============
	12:49
	[CRITICAL]: Resources->
	check this out
	namespace: armo-system, kind: Deployment, name: armo-notification-service
	status:passed()
	Control: C-0014 Access Kubernetes dashboard
	Control: C-0057 Privileged container
	Control: C-0045 Writable hostPath mount
	Control: C-0020 Mount service principal
	Control: C-0012 Applications credentials in configuration files
	Control: C-0048 HostPath mount
	original yaml:
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	 annotations:
	   deployment.kubernetes.io/revision: "8"
	   meta.helm.sh/release-name: armo
	   meta.helm.sh/release-namespace: armo-system
	 creationTimestamp: "2022-01-17T08:42:45Z"
	 generation: 8
	 labels:
	   app: armo-notification-service
	   app.kubernetes.io/managed-by: Helm
	   helm.sh/chart: armo-cluster-components-1.6.7
	   tier: armo-system-control-plane
	 name: armo-notification-service
	 namespace: armo-system
	 resourceVersion: "112927117"
	 uid: 042acd08-1bff-40d8-a6d9-50e388558bcb
	spec:
	 progressDeadlineSeconds: 600
	 replicas: 1
	 revisionHistoryLimit: 10
	 selector:
	   matchLabels:
		 app.kubernetes.io/instance: armo
		 app.kubernetes.io/name: armo-notification-service
		 tier: armo-system-control-plane
	 strategy:
	   rollingUpdate:
		 maxSurge: 25%
		 maxUnavailable: 25%
	   type: RollingUpdate
	 template:
	   metadata:
		 creationTimestamp: null
		 labels:
		   app: armo-notification-service
		   app.kubernetes.io/instance: armo
		   app.kubernetes.io/name: armo-notification-service
		   helm.sh/chart: armo-cluster-components-1.6.7
		   helm.sh/revision: "48"
		   tier: armo-system-control-plane
	   spec:
		 automountServiceAccountToken: false
		 containers:
		 - args:
		   - -alsologtostderr
		   - -v=4
		   - 2>&1
		   env:
		   - name: CA_CUSTOMER_GUID
			 value: XXXXXX
			 valueFrom:
			   configMapKeyRef:
				 key: accountGuid
				 name: armo-be-config
		   - name: CA_CLUSTER_NAME
			 value: XXXXXX
			 valueFrom:
			   configMapKeyRef:
				 key: clusterName
				 name: armo-be-config
		   - name: MASTER_NOTIFICATION_SERVER_HOST
			 value: XXXXXX
			 valueFrom:
			   configMapKeyRef:
				 key: masterNotificationServer
				 name: armo-be-config
		   - name: MASTER_NOTIFICATION_SERVER_ATTRIBUTES
			 value: XXXXXX
		   - name: CA_NOTIFICATION_SERVER_WS_PORT
			 value: XXXXXX
		   - name: CA_NOTIFICATION_SERVER_PORT
			 value: XXXXXX
		   image: quay.io/armosec/notification-server-ubi:89
		   imagePullPolicy: Always
		   name: armo-notification-service
		   ports:
		   - containerPort: 8001
			 name: websocket
			 protocol: TCP
		   - containerPort: 8002
			 name: rest-api
			 protocol: TCP
		   resources:
			 limits:
			   cpu: 100m
			   memory: 50Mi
			 requests:
			   cpu: 10m
			   memory: 10Mi
		   terminationMessagePath: /dev/termination-log
		   terminationMessagePolicy: File
		 dnsPolicy: ClusterFirst
		 restartPolicy: Always
		 schedulerName: default-scheduler
		 securityContext: {}
		 terminationGracePeriodSeconds: 30
		 volumes:
		 - configMap:
			 defaultMode: 420
			 items:
			 - key: clusterData
			   path: clusterData.json
			 name: armo-be-config
		   name: armo-be-config
	==============`
	err := SlackSender("a", "CRITICAL", "#liorabuse", "", "test", msg)
	if err != nil {
		t.Errorf("%s", err)
	}
}
