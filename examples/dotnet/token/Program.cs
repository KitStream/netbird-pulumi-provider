using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet token Group",
    });

    var user = new User("example-user", new UserArgs
    {
        Email = "pulumi-dotnet-token-test@example.com",
        Name = "Pulumi Token Test User",
        IsServiceUser = true,
    });

    var res = new PersonalAccessToken("test-token", new PersonalAccessTokenArgs
    {
        Name = "Pulumi DotNet Token",
        ExpirationDays = 30,
        UserId = user.Id,
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
