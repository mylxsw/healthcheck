<template>
  <div id="app">
    <b-container fluid>
        <b-navbar type="dark" toggleable="md" variant="dark" class="mb-3" sticky>
            <b-navbar-brand href="/">Healthcheck <a href="https://github.com/mylxsw/healthcheck" class="text-white" style="font-size: 30%">{{ version }}</a></b-navbar-brand>
            <b-collapse is-nav id="nav_dropdown_collapse">
                <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex"></ul>
                <b-navbar-nav>
                    <b-nav-item href="/#/" exact>Alerts</b-nav-item>
                    <b-nav-item href="/#/healthchecks" exact>Healthchecks</b-nav-item>
                </b-navbar-nav>
            </b-collapse>
        </b-navbar>
        <div class="main-view">
            <router-view/>
        </div>
    </b-container>
    
  </div>
</template>

<script>
    export default {
        data() {
            return {
                version: 'v-0',
            }
        },
        methods: {
        },
        mounted() {
        },
        beforeMount() {
        }
    }
</script>

<style>
    .container-fluid {
        padding: 0;
    }

    .main-view {
        padding: 15px;
    }

    .th-column-width-limit {
        max-width: 300px;
    }

    @media screen and (max-width: 1366px) {
        .th-autohide-md {
            display: none;
        }
    }
    @media screen and (max-width: 768px) {
        .th-autohide-sm {
            display: none;
        }
        .search-box {
            display: none;
        }
    }

</style>
