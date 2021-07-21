<template>
    <span v-b-tooltip.hover :title="title === '' ? value : title">{{ formatted(value) }}</span>
</template>

<script>
    import moment from 'moment';
    export default {
        name: 'DateTime',
        props: {
            value: String,
            title: String,
        },
        methods: {
            formatted(t) {
                if (t === '0001-01-01T00:00:00Z') {
                    return '-';
                }
                return moment(t).format('YYYY-MM-DD HH:mm:ss');
            }
        }
    }
</script>

