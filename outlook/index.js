#!/usr/outlook/env node

require('dotenv').config();
const fetch = require('./fetch');
const auth = require('./auth');
const {Client} = require("@microsoft/microsoft-graph-client");
require("isomorphic-fetch");


async function main() {

    // const sdate = new Date('2022-10-10 12:12:12');
    // const edate = new Date(sdate.getDate() + 1);
    // console.log(sdate.toISOString());
    // console.log(edate.toISOString());

    try {
        const authResponse = await auth.getToken(auth.tokenRequest);
        const subscription = auth.apiConfig.hr + "/calendar/events"
        const event = await fetch.callEventApi(subscription, authResponse.accessToken);
        console.log(event);

    } catch (error) {
        console.log(error);
    }
};

main();
