
CELLS = Trie.new
DRUGS = Trie.new
TARGETS = Trie.new
TISSUES = Trie.new
DATASETS = Trie.new
GENES = Trie.new

all_cells = Cell.pluck.pluck(2, 0)
all_cells_synonyms = SourceCellName.pluck.pluck(3, 1)

all_cells = all_cells + all_cells_synonyms

all_cells.each do |c|
    unless CELLS.has_key?(c[0].upcase)
        CELLS.add c[0].upcase, c[1]
    end
end


all_drugs = Drug.pluck.pluck(1,0)
all_drugs_synonyms = SourceDrugName.pluck.pluck(3, 1)

all_drugs = all_drugs + all_drugs_synonyms

all_drugs.each do |d|
    unless DRUGS.has_key?(d[0].upcase)
        DRUGS.add(d[0].upcase, d[1])
    end
end

all_target = Target.pluck.pluck(1,0)

all_target.each do |d|
    unless TARGETS.has_key?(d[0].upcase)
        TARGETS.add(d[0].upcase, d[1])
    end
end

all_gene = Gene.pluck.pluck(1,0)


all_gene.each do |g|
    unless GENES.has_key?(g[0].upcase)
        GENES.add(g[0].upcase, g[1])
    end
end

all_tissues = Tissue.pluck.pluck(1,0)

all_tissues.each do |d|
    unless TISSUES.has_key?(d[0].upcase)
        TISSUES.add(d[0].upcase, d[1])
    end
end

all_datasets = Dataset.pluck.pluck(1,0)

all_datasets.each do |d|
    unless DATASETS.has_key?(d[0].upcase)
        DATASETS.add(d[0].upcase, d[1])
    end
end

