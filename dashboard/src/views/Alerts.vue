<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" no-body v-if="alerts.length > 0">
                <b-table :items="alerts" :fields="alert_fields">
                    <template v-slot:cell(id)="row">
                        <b v-b-tooltip.hover title="Name">{{ row.item.healthcheck.name }}</b><br />
                        <i style="color: #808890;" v-b-tooltip.hover title="ID">{{ row.item.healthcheck.id }}</i>
                    </template>
                    <template v-slot:cell(last_check_time)="row">
                        <date-time :value="row.item.last_alive_time" title="Last Check Time"></date-time>
                    </template>
                    <template v-slot:cell(check_type)="row">
                        <b-badge :variant="typeBadge(row.item.healthcheck.check_type)">{{ row.item.healthcheck.check_type }}</b-badge>
                    </template>
                    <template v-slot:cell(status)="row">
                        <b-badge class="mr-2" variant="success" v-if="row.item.alert_times === 0">OK</b-badge>
                        <b-badge class="mr-2" variant="danger"  v-if="row.item.alert_times > 0" @click="row.toggleDetails" style="cursor: pointer">FAIL</b-badge>
                        <br/>
                        <date-time :value="row.item.last_failure_time" title="Last Failure Time" v-if="row.item.alert_times > 0"></date-time>
                        <date-time :value="row.item.last_success_time" title="Last Recovery Time" v-if="row.item.alert_times === 0"></date-time>
                    </template>

                    <template #row-details="row">
                        <b-card>
                            <b-row class="mb-2 pl-3 pr-3">
                                <pre class="text-danger">{{ row.item.last_failure }}</pre>
                            </b-row>
                            <b-button size="sm" @click="row.toggleDetails">Hide</b-button>
                        </b-card>
                    </template>
                </b-table>
            </b-card>
            <b-card class="mb-2" v-if="alerts.length == 0">No Data</b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';
import DateTime from '../components/DateTime.vue';

export default {
        name: 'Alerts',
        components: {DateTime},
        data() {
            return {
                alert_fields: [
                    {key: 'id', label: 'ID/Name'},
                    {key: 'last_check_time', label: 'Last Check Time'},
                    {key: 'check_type', label: 'Check Type'},
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
            typeBadge(typ) {
                switch (typ) {
                    case "http": return "info";
                    case "ping": return "primary";
                }

                return "";
            },
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