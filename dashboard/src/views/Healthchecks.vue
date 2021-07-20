<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" no-body v-if="healthchecks.length > 0">
                <b-table :items="healthchecks" :fields="healthchecks_fields">
                    <template v-slot:cell(last_active_time)="row">
                        <date-time :value="row.item.last_active_time"></date-time>
                    </template>
                    <template v-slot:cell(check_interval)="row">
                        <b class="text-success">{{ row.item.healthcheck.check_interval }}</b> / <b class="text-warning">{{ row.item.healthcheck.loss_threshold }}</b>
                    </template>
                    <template v-slot:cell(id)="row">
                        <b v-b-tooltip.hover title="Name">{{ row.item.healthcheck.name }}</b><br />
                        <i style="color: #808890;" v-b-tooltip.hover title="ID">{{ row.item.healthcheck.id }}</i>
                    </template>
                    <template v-slot:cell(check_type)="row">
                        <b-badge variant="info" @click="row.toggleDetails" style="cursor: pointer" v-b-tooltip.hover title="Show Details">{{ row.item.healthcheck.check_type }}</b-badge>
                    </template>
                    <template v-slot:cell(tags)="row">
                        <div :key="index" v-for="(tag, index) in row.item.healthcheck.tags">
                            <b-badge variant="primary">{{ tag }}</b-badge>
                        </div>
                    </template>
                    <template #row-details="row">
                        <b-card>
                            <b-row class="mb-2 pl-3 pr-3">
                                <pre v-if="row.item.healthcheck.check_type === 'http'">{{ row.item.healthcheck.http }}</pre>
                            </b-row>
                            <b-button size="sm" @click="row.toggleDetails">Hide</b-button>
                        </b-card>
                    </template>
                </b-table>
            </b-card>
            <b-card class="mb-2" v-if="healthchecks.length == 0">No Data</b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';
import DateTime from '../components/DateTime.vue';

export default {
        name: 'Healthchecks',
        components: {DateTime},
        data() {
            return {
                healthchecks_fields: [
                    {key: 'id', label: 'ID/Name'},
                    {key: 'tags', label: 'Tags'},
                    {key: 'check_interval', label: 'Check Interval/Threshold'},
                    {key: 'check_type', label: 'Check Type'},
                    {key: 'last_active_time', label: 'Last Active Time'},
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