using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var res = new Group("test-group", new GroupArgs
    {
        Name = "Pulumi DotNet Group",
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
