{
    Auth: true,
    Rule: [
        {
            URL: [
                '/manager.Manager/Create',
                '/manager.Manager/RemoveID',
                '/manager.Manager/RemoveAccess',
                '/provider.Provider/Create',
                '/provider.Provider/RemoveID',
                '/provider.Provider/RemoveAccess',
            ],
            Bearer: [
                'access token for internal login server',
                'access token for other internal login server',
            ],
        },
        {
            URL: [
                '/provider.Provider/Put',
                '/provider.Provider/Get',
                '/provider.Provider/Keys',
                '/provider.Provider/RemoveKeys',
            ],
            Bearer: [
                'access token for internal server',
                'access token for other internal server',
            ],
        },
        {
            URL: [
                '/logger.Logger/Level',
                '/logger.Logger/File',
                '/logger.Logger/Console',
            ],
            Bearer: [
                'access token for log management',
            ],
        },
    ]
}
