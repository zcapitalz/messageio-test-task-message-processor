#!/bin/bash

ROOT=/home/zcapital/repo/test-tasks/messaggio-message-processor
CERTS=${ROOT}/certs/kafka
CONFIGS=${ROOT}/deploy/kafka

mkdir -p $CERTS

# Generate CA key
openssl req -new -nodes \
   -x509 \
   -days 365 \
   -newkey rsa:2048 \
   -keyout ${CERTS}/ca.key \
   -out ${CERTS}/ca.crt \
   -config ${CONFIGS}/ca.cnf

cat ${CERTS}/ca.crt ${CERTS}/ca.key > ${CERTS}/ca.pem

# Store CA in truststore
keytool \
	-keystore ${CERTS}/kafka.truststore.pkcs12 \
	-alias CARoot \
	-import \
	-file ${CERTS}/ca.crt \
	-storepass confluent  \
	-noprompt \
	-storetype PKCS12
echo "confluent" > ${CERTS}/truststore_creds

for ENTITY in broker-0
do
	# Create certificate
	openssl req -new \
		-newkey rsa:2048 \
		-keyout ${CERTS}/${ENTITY}.key \
		-out ${CERTS}/${ENTITY}.csr \
		-config ${CONFIGS}/${ENTITY}.cnf \
		-nodes

	# Sign certificate
	openssl x509 -req \
		-days 3650 \
		-in ${CERTS}/${ENTITY}.csr \
		-CA ${CERTS}/ca.crt \
		-CAkey ${CERTS}/ca.key \
		-CAcreateserial \
		-out ${CERTS}/${ENTITY}.crt \
		-extfile ${CONFIGS}/${ENTITY}.cnf \
		-extensions v3_req

	# Convert to pkcs12
	openssl pkcs12 -export \
		-in ${CERTS}/${ENTITY}.crt \
		-inkey ${CERTS}/${ENTITY}.key \
		-chain \
		-CAfile ${CERTS}/ca.pem \
		-name ${ENTITY} \
		-out ${CERTS}/${ENTITY}.p12 \
		-password pass:confluent

	# Store in keystore
	keytool -importkeystore \
		-deststorepass confluent \
		-destkeystore ${CERTS}/kafka.${ENTITY}.keystore.pkcs12 \
		-srckeystore ${CERTS}/${ENTITY}.p12 \
		-deststoretype PKCS12  \
		-srcstoretype PKCS12 \
		-noprompt \
		-srcstorepass confluent
	
	# Convert to pem
	openssl pkcs12 -in ${CERTS}/${ENTITY}.p12 -nokeys -out ${CERTS}/${ENTITY}.cer.pem
	openssl pkcs12 -in ${CERTS}/${ENTITY}.p12 -nodes -nocerts -out ${CERTS}/${ENTITY}.key.pem

	# Store credentials
	echo "confluent" > ${CERTS}/${ENTITY}_sslkey_creds
	echo "confluent" > ${CERTS}/${ENTITY}_keystore_creds
done