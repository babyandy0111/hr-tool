const axios = require('axios');

async function callApi(endpoint, accessToken) {
    const options = {
        headers: {
            Authorization: `Bearer ${accessToken}`
        }
    };

    console.log('request made to web API at: ' + new Date().toString());

    try {
        const response = await axios.default.get(endpoint, options);
        return response.data;
    } catch (error) {
        console.log(error)
        return error;
    }
};

async function callEventApi(endpoint, accessToken) {
    const options = {
        headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json'
        }
    }

    const data = {
        "subject": "GAIA Leave List",
        "start": {
            "dateTime": "2022-10-06T10:22:17.532Z",
            "timeZone": "UTC"
        },
        "end": {
            "dateTime": "2022-10-13T10:22:17.532Z",
            "timeZone": "UTC"
        },
        "body": {
            "contentType": "HTML",
            "content": "這邊要放一些請假人名單～～～～～～"
        },
        "attendees": [
            {
                "emailAddress": {
                    "address": "team@gaia.net",
                    "name": "MSG_GAIA CORP"
                },
                "type": "required"
            }
        ],
    }

    try {
        const response = await axios.default.post(endpoint, data, options)
        return response.data;
    } catch (error) {
        // console.log(error)
        return error;
    }
};

module.exports = {
    callApi: callApi,
    callEventApi: callEventApi
};
