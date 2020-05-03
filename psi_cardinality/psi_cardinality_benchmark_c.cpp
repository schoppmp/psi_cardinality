#include "absl/strings/str_cat.h"
#include "benchmark/benchmark.h"
#include "psi_cardinality/psi_cardinality_client_c.h"
#include "psi_cardinality/psi_cardinality_server_c.h"

namespace psi_cardinality {
namespace {

void BM_ServerBindingSetup(benchmark::State& state, double fpr) {
    auto server_ = psi_cardinality_server_create_with_new_key();

    int num_inputs = state.range(0);
    int num_client_inputs = 10000;
    std::vector<std::string> orig_inputs(num_inputs);
    std::vector<server_buffer_t> inputs(num_inputs);

    for (int i = 0; i < num_inputs; i++) {
        orig_inputs[i] = absl::StrCat("Element", i);
        inputs[i] = {orig_inputs[i].c_str(), orig_inputs[i].size()};
    }
    std::string setup;
    int64_t elements_processed = 0;
    for (auto _ : state) {
        char* server_setup = nullptr;
        size_t server_setup_buff_len = 0;
        psi_cardinality_server_create_setup_message(
            server_, fpr, num_client_inputs, inputs.data(), inputs.size(),
            &server_setup, &server_setup_buff_len);

        ::benchmark::DoNotOptimize(
            std::string(server_setup, server_setup_buff_len));
        elements_processed += num_inputs;
        psi_cardinality_server_delete_buffer(server_, &server_setup);
    }
    state.counters["SetupSize"] = benchmark::Counter(
        static_cast<double>(setup.size()), benchmark::Counter::kDefaults,
        benchmark::Counter::kIs1024);
    state.counters["ElementsProcessed"] = benchmark::Counter(
        static_cast<double>(elements_processed), benchmark::Counter::kIsRate);
}

// Range is for the number of inputs, and the captured argument is the false
// positive rate for 10k client queries.
BENCHMARK_CAPTURE(BM_ServerBindingSetup, fpr = 0.001, 0.001)
    ->RangeMultiplier(100)
    ->Range(1, 1000000);
BENCHMARK_CAPTURE(BM_ServerBindingSetup, fpr = 0.000001, 0.000001)
    ->RangeMultiplier(100)
    ->Range(1, 1000000);

void BM_ClientBindingCreateRequest(benchmark::State& state) {
    auto client_ = psi_cardinality_client_create();

    int num_inputs = state.range(0);
    std::vector<std::string> inputs_orig(num_inputs);
    std::vector<client_buffer_t> inputs(num_inputs);

    for (int i = 0; i < num_inputs; i++) {
        inputs_orig[i] = absl::StrCat("Element", i);
        inputs[i] = {inputs_orig[i].c_str(), inputs_orig[i].size()};
    }
    std::string request;
    int64_t elements_processed = 0;
    for (auto _ : state) {
        char* client_request = {0};
        size_t req_len = 0;
        psi_cardinality_client_create_request(
            client_, inputs.data(), inputs.size(), &client_request, &req_len);

        ::benchmark::DoNotOptimize(std::string(client_request, req_len));
        elements_processed += num_inputs;

        psi_cardinality_client_delete_buffer(client_, &client_request);
    }
    state.counters["RequestSize"] = benchmark::Counter(
        static_cast<double>(request.size()), benchmark::Counter::kDefaults,
        benchmark::Counter::kIs1024);
    state.counters["ElementsProcessed"] = benchmark::Counter(
        static_cast<double>(elements_processed), benchmark::Counter::kIsRate);
}
// Range is for the number of inputs.
BENCHMARK(BM_ClientBindingCreateRequest)->RangeMultiplier(10)->Range(1, 10000);

void BM_ServerBindingProcessRequest(benchmark::State& state) {
    auto client_ = psi_cardinality_client_create();
    auto server_ = psi_cardinality_server_create_with_new_key();

    int num_inputs = state.range(0);
    std::vector<std::string> inputs_orig(num_inputs);
    std::vector<client_buffer_t> inputs(num_inputs);
    for (int i = 0; i < num_inputs; i++) {
        inputs_orig[i] = absl::StrCat("Element", i);
        inputs[i] = {inputs_orig[i].c_str(), inputs_orig[i].size()};
    }

    char* client_request = nullptr;
    size_t req_len = 0;
    psi_cardinality_client_create_request(client_, inputs.data(), inputs.size(),
                                          &client_request, &req_len);

    std::string response;

    int64_t elements_processed = 0;
    for (auto _ : state) {
        char* server_response = nullptr;
        size_t response_len = 0;
        psi_cardinality_server_process_request(server_,
                                               {client_request, req_len},
                                               &server_response, &response_len);
        response = std::string(server_response, response_len);
        ::benchmark::DoNotOptimize(response);
        elements_processed += num_inputs;

        psi_cardinality_server_delete_buffer(server_, &server_response);
    }
    state.counters["ResponseSize"] = benchmark::Counter(
        static_cast<double>(response.size()), benchmark::Counter::kDefaults,
        benchmark::Counter::kIs1024);
    state.counters["ElementsProcessed"] = benchmark::Counter(
        static_cast<double>(elements_processed), benchmark::Counter::kIsRate);

    psi_cardinality_client_delete_buffer(client_, &client_request);
}

// Range is for the number of inputs.
BENCHMARK(BM_ServerBindingProcessRequest)->RangeMultiplier(10)->Range(1, 10000);

void BM_ClientBindingProcessResponse(benchmark::State& state) {
    auto client_ = psi_cardinality_client_create();
    auto server_ = psi_cardinality_server_create_with_new_key();

    int num_inputs = state.range(0);
    double fpr = 1. / (1000000);
    std::vector<std::string> inputs_orig(num_inputs);
    std::vector<server_buffer_t> srv_inputs(num_inputs);
    std::vector<client_buffer_t> cl_inputs(num_inputs);
    for (int i = 0; i < num_inputs; i++) {
        inputs_orig[i] = absl::StrCat("Element", i);
        srv_inputs[i] = {inputs_orig[i].c_str(), inputs_orig[i].size()};
        cl_inputs[i] = {inputs_orig[i].c_str(), inputs_orig[i].size()};
    }

    char* server_setup = nullptr;
    size_t server_setup_buff_len = 0;
    psi_cardinality_server_create_setup_message(
        server_, fpr, num_inputs, srv_inputs.data(), srv_inputs.size(),
        &server_setup, &server_setup_buff_len);

    char* client_request = nullptr;
    size_t req_len = 0;
    psi_cardinality_client_create_request(
        client_, cl_inputs.data(), cl_inputs.size(), &client_request, &req_len);

    char* server_response = nullptr;
    size_t response_len = 0;
    psi_cardinality_server_process_request(server_, {client_request, req_len},
                                           &server_response, &response_len);

    int64_t elements_processed = 0;
    for (auto _ : state) {
        int64_t count = psi_cardinality_client_process_response(
            client_, server_setup, server_response);

        ::benchmark::DoNotOptimize(count);
        elements_processed += num_inputs;
    }

    psi_cardinality_server_delete_buffer(server_, &server_setup);
    psi_cardinality_server_delete_buffer(server_, &server_response);
    psi_cardinality_client_delete_buffer(client_, &client_request);

    state.counters["ElementsProcessed"] = benchmark::Counter(
        static_cast<double>(elements_processed), benchmark::Counter::kIsRate);
}
// Range is for the number of inputs.
BENCHMARK(BM_ClientBindingProcessResponse)
    ->RangeMultiplier(10)
    ->Range(1, 10000);

}  // namespace
}  // namespace psi_cardinality
