<script>
export default {
    data: function() {
        return {
            username: null,
            password: null
        }
    },
    methods: {
        signIn: function() {
            console.log(this.username);
        }
    }
}
</script>