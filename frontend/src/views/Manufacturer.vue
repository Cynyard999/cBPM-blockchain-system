<template>
    <div>
        <h1>Manufacturer</h1>
    </div>
</template>

<script>
    import {request} from '../api/axios'

    export default {
        mounted(): void {
            let params = {
                orgType: "supplier",
                channelName: "mischannel",
                contractName: "mischaincode",
                function: "CreateAsset",
                transient: {
                    asset: {
                        assetName: "asset1",
                        assetPrice: 100.1,
                        shippingAddress: "PlaceB",
                        publicDescription: "good"
                    }
                }
            };
            request('/invoke', params, 'post').then(response => {
                console.log(response)
                if (response.status === 200) {
                    console.log(response.data)
                }
            })
                .catch(function (error) {
                    console.log(error);
                })
        }
    }
</script>

<style scoped>

</style>