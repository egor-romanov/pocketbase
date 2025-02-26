<script>
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";

    export let collection = new Collection();

    let responseTab = 200;
    let sdkTab = "JavaScript";
    let responses = [];
    let sdkExamples = [];

    $: adminsOnly = collection?.createRule === null;

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(CommonHelper.dummyCollectionRecord(collection), null, 2),
        },
        {
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${collection?.schema?.[0]?.name}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
        {
            code: 403,
            body: `
                {
                  "code": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `,
        },
    ];

    $: sdkExamples = [
        {
            lang: "JavaScript",
            code: `
                import PocketBase from 'pocketbase';

                const client = new PocketBase("${ApiClient.baseUrl}");

                const data = { ... };

                client.Records.create("${collection?.name}", data)
                    .then(function (record) {
                        // success...
                    }).catch(function (error) {
                        // error...
                    });
            `,
        },
    ];
</script>

<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">Requires <code>Authorization: Admin TOKEN</code> header</p>
    {/if}
</div>

<div class="content m-b-base">
    <p>Create a new <strong>{collection.name}</strong> record.</p>
    <p>
        Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.
    </p>
    <p>
        File upload is supported only via <code>multipart/form-data</code>.
    </p>
</div>

<div class="section-title">Client SDKs example</div>
<div class="tabs m-b-lg">
    <div class="tabs-header compact left">
        {#each sdkExamples as example (example.lang)}
            <button
                class="tab-item"
                class:active={sdkTab === example.lang}
                on:click={() => (sdkTab = example.lang)}
            >
                {example.lang}
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each sdkExamples as example (example.lang)}
            <div class="tab-item" class:active={sdkTab === example.lang}>
                <CodeBlock content={example.code} />
            </div>
        {/each}
    </div>
</div>

<div class="section-title">Body Parameters</div>
<table class="table-compact table-border m-b-lg">
    <thead>
        <tr>
            <th>Param</th>
            <th>Type</th>
            <th width="50%">Description</th>
        </tr>
    </thead>
    <tbody>
        {#each collection?.schema as field (field.name)}
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if field.required}
                            <span class="label label-success">Required</span>
                        {:else}
                            <span class="label label-warning">Optional</span>
                        {/if}
                        <span>{field.name}</span>
                    </div>
                </td>
                <td>
                    <span class="label">{CommonHelper.getFieldValueType(field)}</span>
                </td>
                <td>
                    {#if field.type === "text"}
                        Plain text value.
                    {:else if field.type === "number"}
                        Number value.
                    {:else if field.type === "json"}
                        JSON array or object.
                    {:else if field.type === "email"}
                        Email address.
                    {:else if field.type === "url"}
                        URL address.
                    {:else if field.type === "file"}
                        FormData object.<br />
                        Set to <code>null</code> to delete already uploaded file(s).
                    {:else if field.type === "relation"}
                        Relation record {field.options?.maxSelect > 1 ? "ids" : "id"}.
                    {:else if field.type === "user"}
                        User {field.options?.maxSelect > 1 ? "ids" : "id"}.
                    {/if}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<div class="section-title">Responses</div>
<div class="tabs">
    <div class="tabs-header compact left">
        {#each responses as response (response.code)}
            <button
                class="tab-item"
                class:active={responseTab === response.code}
                on:click={() => (responseTab = response.code)}
            >
                {response.code}
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each responses as response (response.code)}
            <div class="tab-item" class:active={responseTab === response.code}>
                <CodeBlock content={response.body} />
            </div>
        {/each}
    </div>
</div>
