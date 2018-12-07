<template>
    <aside class="sidebar sidebar-left">
        <nav class="sidebar-sticky">
            <h6 class="sidebar-heading">Library</h6>
            <ul class="nav flex-column">
                <li v-bind:class="{ active: isActive('documents') }" class="nav-item">
                    <router-link :to="{ name: 'documents' }" v-on:click.native="navigate">
                        <span class="entypo-window"></span>Documents
                    </router-link>
                </li>
            </ul>
            <h6 class="sidebar-heading">Shared</h6>
            <ul class="nav flex-column">
                <li v-bind:class="{ active: isActive('sharedActivity') }" class="nav-item">
                    <router-link :to="{ name: 'sharedActivity' }" v-on:click.native="navigate">
                        <span class="entypo-window"></span>Activity
                    </router-link>
                </li>
            </ul>
        </nav>
    </aside>
</template>

<script>
export default {
    data: function() {
        return {};
    },
    methods: {
        isActive: function(routeName) {
            return this.$route.name === routeName;
        },
        navigate: function() {
        }
    }
}
</script>

<style lang="scss">
.sidebar {
    background-image: linear-gradient(-179deg, #f7f7f6 0%, #ebebeb 100%);
}
.sidebar-left {
    position: relative;
    vertical-align: top;
    width: 300px;
}
@include media-breakpoint-up(md) {
    .sidebar-left {
        bottom: 0;
        left: 0;
        padding-bottom: 50px;
        position: fixed;
        top: 60px;
    }
}
</style>