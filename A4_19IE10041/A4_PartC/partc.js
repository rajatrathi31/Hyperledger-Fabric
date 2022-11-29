/*
Assignment 4
Part C

Name: Rajat Rathi
Roll No.: 19IE10041

*/

// Pre-requisite to run file: npm install prompt-sync

const FabricCAServices = require('fabric-ca-client') 
const { Wallets, Gateway } = require('fabric-network')

const fs = require('fs')
const path = require('path')
const prompt = require("prompt-sync")({ sigint: true });

async function main(){
    var ccpPath = path.resolve('../organizations/peerOrganizations/org1.example.com/connection-org1.json') 
    var ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'))

    var caInfo = ccp.certificateAuthorities['ca.org1.example.com'] 
    var caTLSCACerts = caInfo.tlsCACerts.pem 
    var ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName) 

    const walletPath = path.join(process.cwd(), 'wallet')
    const wallet = await Wallets.newFileSystemWallet(walletPath)

    var adminIdentity = await wallet.get("admin1")
    if (adminIdentity){
        console.log('An identity for the admin user for organization 1 already exists in the wallet')
    }
    else {
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' }); 
        const x509Identity = { 
            credentials: { 
                certificate: enrollment.certificate, 
                privateKey: enrollment.key.toBytes(), 
            }, 
            mspId: 'Org1MSP', 
            type: 'X.509',
        }; 

        await wallet.put("admin1", x509Identity) 
        console.log("admin for organization 1 enrolled and saved into wallet successfully")  
        adminIdentity = await wallet.get("admin1")
    }

    var userIdentity1 = await wallet.get('User_1')

    if (userIdentity1){
        console.log('An identity for the client user "User_1" already exists in the wallet')
    }
    else{
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type); 
        const adminUser = await provider.getUserContext(adminIdentity, 'admin'); 

        const secret = await ca.register({
            affiliation: 'org1.department1', 
            enrollmentID: 'User_1', 
            role: 'client'}, adminUser); 
        const enrollment = await ca.enroll({
            enrollmentID: 'User_1',
            enrollmentSecret: secret}); 
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes()},
            mspId: 'Org1MSP', 
            type: 'X.509', }; 
        await wallet.put('User_1', x509Identity) 
        console.log("Enrolled User_1 and saved to wallet");
        userIdentity1 = await wallet.get('User_1')
    }

    ccpPath = path.resolve('../organizations/peerOrganizations/org2.example.com/connection-org2.json') 
    ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'))

    caInfo = ccp.certificateAuthorities['ca.org2.example.com'] 
    caTLSCACerts = caInfo.tlsCACerts.pem 
    ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName) 

    adminIdentity = await wallet.get("admin2")
    if (adminIdentity){
        console.log('An identity for the admin user for organization 2 already exists in the wallet')
    }
    else {
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' }); 
        const x509Identity = { 
            credentials: { 
                certificate: enrollment.certificate, 
                privateKey: enrollment.key.toBytes(), 
            }, 
            mspId: 'Org2MSP', 
            type: 'X.509',
        }; 

        await wallet.put("admin2", x509Identity) 
        console.log("admin for organization 2 enrolled and saved into wallet successfully")  
        adminIdentity = await wallet.get("admin2")
    }

    var userIdentity2 = await wallet.get('User_2')

    if (userIdentity2){
        console.log('An identity for the client user "User_2" already exists in the wallet')
    }
    else{
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type); 
        const adminUser = await provider.getUserContext(adminIdentity, 'admin'); 

        const secret = await ca.register({
            affiliation: 'org2.department1', 
            enrollmentID: 'User_2', 
            role: 'client'}, adminUser); 
        const enrollment= await ca.enroll({
            enrollmentID: 'User_2',
            enrollmentSecret: secret}); 
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes()},
            mspId: 'Org2MSP', 
            type: 'X.509', }; 
        await wallet.put('User_2', x509Identity) 
        console.log('Enrolled "User_2" and saved to wallet');
        userIdentity2 = await wallet.get('User_2')
    }

    const gateway1 = new Gateway(); 
    await gateway1.connect(ccp, {wallet, identity:'User_1', discovery: {enabled: true, asLocalhost: true}}) 
    const network1 = await gateway1.getNetwork('mychannel') // select the contract 
    const contract1 = network1.getContract("basic")

    const gateway2 = new Gateway(); 
    await gateway2.connect(ccp, {wallet, identity:'User_2', discovery: {enabled: true, asLocalhost: true}}) 
    const network2 = await gateway1.getNetwork('mychannel') // select the contract 
    const contract2 = network2.getContract("basic")
    
    var org_num = 1;
    while (true) {
        var instruction = prompt("Enter instruction(EXIT to exit) : ");
        if (instruction == "INSERT") {
            var num = prompt("Enter the number to insert : ");
            if (org_num == 1) {
                err = await contract1.submitTransaction("Insert", num)
                if (err.toString() != "") {      
                    console.log(err);
				} else {
                    org_num = 2
					console.log("The number has been inserted in the tree.");
				}
            } else {
                err = await contract2.submitTransaction("Insert", num)
                if (err.toString() != "") {     
                    console.log(err);
				} else {
                    org_num = 1
					console.log("The number has been inserted in the tree.");
				}
            }
        } else if (instruction == "PREORDER") {
            if (org_num == 1) {
                result = await contract1.evaluateTransaction("Preorder") 
                org_num = 2
				console.log("Preorder Traversal : ", result.toString()) 
            } else {
                result = await contract2.evaluateTransaction("Preorder")
				org_num = 1 
                console.log("Preorder Traversal : ", result.toString()) 
            }
        } else if (instruction == "INORDER") {
            if (org_num == 1) {
                result = await contract1.evaluateTransaction("Inorder") 
				org_num = 2
                console.log("Inorder Traversal : " + result.toString()) 
            } else {
                result = await contract2.evaluateTransaction("Inorder") 
				org_num = 1
                console.log("Inorder Traversal : " + result.toString()) 
            }
        } else if (instruction == "TREEHEIGHT") {
            if (org_num == 1) {
                result = await contract1.evaluateTransaction("TreeHeight")
				org_num = 2 
                console.log("Tree Height : ", result.toString()) 
            } else {
                result= await contract2.evaluateTransaction("TreeHeight")
				org_num = 1 
                console.log("Tree Height : ", result.toString()) 
            }
        } else if (instruction == "DELETE") {
            var num = prompt("Enter the number to delete : ");
            if (org_num == 1) {
                err = await contract1.submitTransaction("Delete",num)
                if (err.toString() != "") {     
                    console.log(err);
				} else {
                    org_num = 2
					console.log("The number has been deleted from the tree.");
				}
            } else {
                err = await contract2.submitTransaction("Delete",num)
                if (err.toString() != "") {     
                    console.log(err);
				} else {
                    org_num = 1
					console.log("The number has been deleted from the tree.");
				}
            }
        } else if (instruction == "EXIT") {
            break;
		} else {
            console.log("Enter valid instruction (Capital Letters only)")
		}
    }
    gateway1.disconnect()
    gateway2.disconnect()
}

main();