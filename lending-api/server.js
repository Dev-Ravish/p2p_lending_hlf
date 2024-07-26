const express = require('express');
const bodyParser = require('body-parser');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');
const mysql = require('mysql');

const app = express();
const port = 3000;

app.use(bodyParser.json());

// Database connection
const connection = mysql.createConnection({
    host: 'localhost',
    port: 3306,
    user: 'mariadbuser',
    password: 'yourPassword',
    database: 'mariadbdatabase'
});

// Connect to Fabric network
async function connectToNetwork() {
    const ccpPath = path.resolve(__dirname, '..','test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'appUser',
        discovery: { enabled: true, asLocalhost: true }
    });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('basic');
    return contract;
}

//Endpoint to create the new user
app.post('/createUser', async (req, res) => {
    try {
        const { id, name, risk, fund, email_id, account } = req.body;
        const contract = await connectToNetwork();
        await contract.submitTransaction('CreateAccount', id, name, risk, fund);

        // Insert user details into MariaDB
        const insertUserQuery = 'INSERT INTO users (name, email_id, account, user_id) VALUES (?, ?, ?, ?)';
        connection.query(insertUserQuery, [name, email_id, account, id], (error, results) => {
            if (error) {
                return res.status(500).send(`Failed to insert user into database: ${error}`);
            }
            res.status(200).send(`Transaction has been submitted: ${id} with name ${name} created and added to the database`);
        });
    } catch (error) {
        res.status(500).send(`Failed to submit transaction: ${error}`);
    }
});

app.get('/getAllUsers', (req, res) => {
    const query = 'SELECT * FROM users';

    connection.query(query, (error, results) => {
        if (error) {
            return res.status(500).send(`Failed to retrieve data from database: ${error}`);
        }
        res.status(200).json(results);
    });
});

//Endpoint to call 'borrow' function
app.put('/borrow', async (req, res) => {
    try {
        const { borrowerID, fundsNeeded } = req.body
        const contract = await connectToNetwork()
        await contract.submitTransaction('Borrow', borrowerID, fundsNeeded)
        res.status(200).send(`Transaction has been submitted: ${borrowerID} borrowed ${fundsNeeded}`)
    } catch (error) {
        res.status(500).send(`Failed to submit transaction: ${error}`)
    }
});

//Endpoint to add funds
app.put('/addFunds', async(req, res)=>{
    try {
        const {accountId, amount} = req.body
        const contract = await connectToNetwork()
        await contract.submitTransaction('AddFunds', accountId, amount)
        res.status(200).send(`Transaction has been submitted: ${accountId} added ${amount} funds`)
    }catch(err){
        res.status(500).send(`Failed to submit transaction: ${err}`)
    }
})

//Endpoint to give loans
app.put('/giveLoan', async(req, res)=>{
    try {
        const {borrowerId, lenderId, tid} = req.body
        const contract = await connectToNetwork()
        await contract.submitTransaction('FundGiving', lenderId, borrowerId, tid)
        res.status(200).send(`Transaction has been submitted: ${lenderId} gave loan to ${borrowerId}`)
    }catch(err){
        res.status(500).send(`Failed to submit transaction: ${err}`)
    }
})

//Endpoint to repay loans
app.put('/repayLoan', async(req, res)=>{
    try {
        const {borrowerId, lenderId, tid, amount} = req.body
        const contract = await connectToNetwork()
        await contract.submitTransaction('LoanRepayment', lenderId, borrowerId, tid, amount)
        res.status(200).send(`Transaction has been submitted: ${borrowerId} repaid the required loan to ${lenderId}`)
    }catch(err){
        res.status(500).send(`Failed to submit transaction: ${err}`)
    }
})

//Endpoint to call 'updateRisk' function
app.put('/updateRisk', async (req, res) => {
    try {
        const { accountID, risk, auto } = req.body;
        const contract = await connectToNetwork();
        await contract.submitTransaction('updateRisk', accountID, risk, auto);
        res.status(200).send(`Risk updated for account: ${accountID}`);
    } catch (error) {
        res.status(500).send(`Failed to submit transaction: ${error}`);
    }
});

//Endpoint to initialize network with some fake data
app.get('/initialize', async(req, res) =>{
    try{
        const contract = await connectToNetwork();
        const result = await contract.submitTransaction('InitLedger');
        console.log("result done", result)

        res.status(200).json({"message":"Initializing succesfull!!"});

    }catch(err){
        res.status(500).send(`Failed to evaluate transaction: ${err}`);
    }
})

// Endpoint to query all accounts
app.get('/getAllAccounts', async (req, res) => {
    try {
        console.log("connection start")

        const contract = await connectToNetwork();
        console.log("connection done")
        const result = await contract.evaluateTransaction('GetAllAccounts');
        console.log("result done")

        res.status(200).json(JSON.parse(result.toString()));
    } catch (error) {
        res.status(500).send(`Failed to evaluate transaction: ${error}`);
    }
});


app.get('/getAllTransactions', async (req, res) => {
    try {

        const contract = await connectToNetwork();
        const result = await contract.evaluateTransaction('GetAllTransactions');

        res.status(200).json(JSON.parse(result.toString()));
    } catch (error) {
        res.status(500).send(`Failed to evaluate transaction: ${error}`);
    }
});

app.get('/getAnAccount', async (req, res) => {
    try {
        const {bid, lid} = req.body

        const contract = await connectToNetwork();
        const result = await contract.evaluateTransaction('ReadAccount', borrowerId=bid, lenderId=lid);

        res.status(200).json(JSON.parse(result.toString()));
    } catch (error) {
        res.status(500).send(`Failed to evaluate transaction: ${error}`);
    }
});

app.listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
});
