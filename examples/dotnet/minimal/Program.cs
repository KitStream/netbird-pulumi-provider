using System.Collections.Generic;
using Pulumi;
using Pulumi.KitStream.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("test-group", new GroupArgs
    {
        Name = "Pulumi DotNet Test Group",
    });

    return new Dictionary<string, object?>
    {
        ["groupName"] = group.Name,
    };
});
