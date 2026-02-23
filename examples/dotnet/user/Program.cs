using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet user Group",
    });

    var res = new User("test-user", new UserArgs
    {
        Email = "pulumi-DotNet@example.com",
        Name = "Pulumi DotNet User",
        IsServiceUser = true,
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Email,
    };
});
