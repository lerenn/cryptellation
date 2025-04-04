import * as hey_api from '@hey-api/client-fetch';

import * as types from './generated/types.gen';
import * as sdk from './generated/sdk.gen';

export class Client {
    private client : hey_api.Client;

    constructor(options : hey_api.ClientOptions) {
        this.client = hey_api.createClient(options);
    }

    async getInfo() : Promise<({
        data: types.SystemInformation;
        error: undefined;
    } | {
        data: undefined;
        error: unknown;
    }) & {
        request: Request;
        response: Response;
    }>{
        return sdk.getInfo({
            client: this.client,
        });
    }
}