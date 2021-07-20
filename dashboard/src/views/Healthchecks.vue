<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" no-body v-if="healthchecks.length > 0">
                <b-table :items="healthchecks" :fields="healthchecks_fields">
                    <template v-slot:cell(status)="row">
                        <b-badge class="mr-2" variant="success" v-if="row.item.alert_times == 0">正常</b-badge>
                        <b-badge class="mr-2" variant="danger"  v-if="row.item.alert_times > 0">失败</b-badge>
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

export default {
        name: 'Shares',
        components: {},
        data() {
            return {
                healthchecks_fields: [
                    {key: 'id', label: 'ID'},
                    {key: 'name', label: 'Name'},
                    {key: 'check_interval', label: 'Check Interval'},
                    {key: 'loss_threshold', label: 'Loss Threshold'},
                    {key: 'check_type', label: 'Check Type'},
                    {key: 'status', label: 'Status'},
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