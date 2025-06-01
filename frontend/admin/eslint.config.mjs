import withNuxt from './.nuxt/eslint.config.mjs';
//import eslintPluginPrettier from "eslint-plugin-prettier";
//import eslintConfigPrettier from 'eslint-config-prettier';

export default withNuxt({
    files: ['**/*.ts', '**/*.vue'],
    rules: {
        'no-unused-vars': 'warn',
        '@typescript-eslint/no-unused-vars': 'warn',
        'vue/no-multiple-template-root': 'warn',
        'vue/multi-word-component-names': 'off',
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/ban-ts-comment': 'off',
    },
});
