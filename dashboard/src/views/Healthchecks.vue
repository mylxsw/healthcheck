<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" no-body v-if="healthchecks.length > 0">
                <b-table :items="healthchecks" :fields="healthchecks_fields">
                    <template v-slot:cell(last_success_time)="row">
                        <date-time :value="row.item.last_success_time"></date-time>
                    </template>
                </b-table>
            </b-card>
            <b-card class="mb-2" v-if="healthchecks.length == 0">当前没有相关健康检查</b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';
import moment from 'moment';
import DateTime from '../components/DateTime.vue';

export default {
        name: 'Healthchecks',
        components: {DateTime},
        data() {
            return {
                healthchecks_fields: [
                    {key: 'healthcheck.id', label: 'ID'},
                    {key: 'healthcheck.name', label: 'Name'},
                    {key: 'healthcheck.check_interval', label: 'Check Interval'},
                    {key: 'healthcheck.loss_threshold', label: 'Loss Threshold'},
                    {key: 'healthcheck.check_type', label: 'Check Type'},
                    {key: 'last_success_time', label: 'Last Success Time'},
                ],
                healthchecks: [],
            };
        },
        computed: {
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            reload() {
                axios.get('/api/healthchecks/').then(resp => {
                    this.healthchecks = resp.data;
                }).catch(error => {this.ToastError(error)});
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>