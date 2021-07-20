<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" no-body v-if="alerts.length > 0">
                <b-table :items="alerts" :fields="alert_fields">
                    <template v-slot:cell(status)="row">
                        <b-badge class="mr-2" variant="success" v-if="row.item.alert_times == 0">正常</b-badge>
                        <b-badge class="mr-2" variant="danger"  v-if="row.item.alert_times > 0">失败</b-badge>
                    </template>
                </b-table>
            </b-card>
            <b-card class="mb-2" v-if="alerts.length == 0">当前没有相关告警</b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';
import moment from 'moment';

export default {
        name: 'Alerts',
        components: {},
        data() {
            return {
                alert_fields: [
                    {key: 'healthcheck.id', label: 'ID'},
                    {key: 'healthcheck.name', label: 'Name'},
                    {key: 'healthcheck.check_interval', label: 'Check Interval'},
                    {key: 'healthcheck.loss_threshold', label: 'Loss Threshold'},
                    {key: 'healthcheck.check_type', label: 'Check Type'},
                    {key: 'status', label: 'Status'},
                ],
                alerts: [],
            };
        },
        computed: {
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            reload() {
                axios.get('/api/alerts/').then(resp => {
                    this.alerts = resp.data;
                }).catch(error => {this.ToastError(error)});
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>