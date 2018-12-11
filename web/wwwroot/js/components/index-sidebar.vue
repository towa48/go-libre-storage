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

    .sidebar-heading {
        font-size: $sidebar-head-font-size;
        color: #707e8c;
        line-height: $sidebar-head-line-height;
        text-shadow: 0 1px 0 rgba(255,255,255,0.46);
        text-transform: uppercase;
        display: inline-block;
        padding: 5px 9px;
    }
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